package annotation

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/spf13/cast"

	"github.com/nats-io/nats.go"

	"github.com/spf13/viper"

	"github.com/nkhang/pluto/internal/dataset"
	"github.com/nkhang/pluto/internal/label"
	"github.com/nkhang/pluto/internal/project"
	"github.com/nkhang/pluto/internal/task"
	"github.com/nkhang/pluto/internal/workspace"
	"github.com/nkhang/pluto/pkg/errors"
	"github.com/nkhang/pluto/pkg/logger"
	"github.com/nkhang/pluto/pkg/util/clock"
)

type Service interface {
	CreateTask(projectID, datasetID uint64, tasks []task.Task) error
	UpdateProject(projectID uint64) error
	UpdateDataset(datasetID uint64) error
	GetLabelCount(projectID, labelID uint64) (LabelStatsObject, error)
	CreateTaskWithNATS(projectID, datasetID uint64, tasks []task.Task) error
	GetImageStats(projectID uint64) (obj LabelStatsObject, err error)
}

type service struct {
	workspaceRepo      workspace.Repository
	projectRepo        project.Repository
	datasetRepo        dataset.Repository
	labelRepo          label.Repository
	client             http.Client
	nc                 *nats.EncodedConn
	annotationBasePath string
}

func NewService(workspaceRepo workspace.Repository,
	projectRepo project.Repository,
	datasetRepo dataset.Repository,
	labelRepo label.Repository) *service {
	client := http.Client{}
	annotationBase := viper.GetString("annotation.baseurl")
	return &service{
		workspaceRepo:      workspaceRepo,
		projectRepo:        projectRepo,
		datasetRepo:        datasetRepo,
		labelRepo:          labelRepo,
		client:             client,
		annotationBasePath: annotationBase,
	}
}

func (s *service) GetLabelCount(projectID, labelID uint64) (obj LabelStatsObject, err error) {
	path := s.annotationBasePath + "/stats/labels"
	u, err := url.Parse(path)
	if err != nil {
		err = errors.AnnotationCannotParseURL.WrapF(err, "cannot parse url %s", path)
		return
	}
	q := u.Query()
	q.Set("project_id", cast.ToString(projectID))
	q.Set("label_id", cast.ToString(labelID))
	u.RawQuery = q.Encode()
	logger.Infof("request: %s", u.String())
	resp, err := s.client.Get(u.String())
	if err != nil {
		err = errors.AnnotationCannotGetFromServer.WrapF(err, "cannot get annotation label statistic from server")
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.AnnotationCannotReadBody.Wrap(err, "cannot get response body")
		return
	}
	logger.Infof("got: %s", b)
	var respObj LabelStatsResponse
	err = json.Unmarshal(b, &respObj)
	if err != nil {
		err = errors.AnnotationCannotReadBody.Wrap(err, "cannot parse json body of response")
		return
	}
	if respObj.Status != 1 {
		err = errors.AnnotationCannotGetFromServer.NewWithMessageF("error getting from annotation server. msg: %s", respObj.Message)
	}
	return respObj.Data, nil
}

func (s *service) GetImageStats(projectID uint64) (obj LabelStatsObject, err error) {
	path := s.annotationBasePath + "/stats/images"
	u, err := url.Parse(path)
	if err != nil {
		err = errors.AnnotationCannotParseURL.WrapF(err, "cannot parse url %s", path)
		return
	}
	q := u.Query()
	q.Set("project_id", cast.ToString(projectID))
	u.RawQuery = q.Encode()
	logger.Infof("[ANNOTATION] - request URL: %s", u.String())
	resp, err := s.client.Get(u.String())
	if err != nil {
		err = errors.AnnotationCannotGetFromServer.WrapF(err, "cannot get annotation label statistic from server")
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.AnnotationCannotReadBody.Wrap(err, "cannot get response body")
		return
	}
	logger.Infof("[ANNOTATION] - response: %s", b)
	var respObj LabelStatsResponse
	err = json.Unmarshal(b, &respObj)
	if err != nil {
		err = errors.AnnotationCannotReadBody.Wrap(err, "cannot parse json body of response")
		return
	}
	if respObj.Status != 1 {
		err = errors.AnnotationCannotGetFromServer.NewWithMessageF("error getting from annotation server. msg: %s", respObj.Message)
	}
	return respObj.Data, nil
}

func (s *service) CreateTaskWithNATS(projectID, datasetID uint64, tasks []task.Task) error {
	p, err := s.projectRepo.Get(projectID)
	if err != nil {
		return err
	}
	message, err := NewBuilder(
		s.workspaceRepo,
		s.projectRepo,
		s.datasetRepo,
		s.labelRepo).
		WithWorkspace(p.WorkspaceID).
		WithProject(projectID).
		WithDataset(datasetID).
		WithTasks(tasks).
		WithLabels(projectID).
		Build()
	if err != nil {
		return err
	}
	return s.pushWithNATS(message)
}

func (s *service) CreateTask(projectID, datasetID uint64, tasks []task.Task) error {
	p, err := s.projectRepo.Get(projectID)
	if err != nil {
		return err
	}
	message, err := NewBuilder(
		s.workspaceRepo,
		s.projectRepo,
		s.datasetRepo,
		s.labelRepo).
		WithWorkspace(p.WorkspaceID).
		WithProject(projectID).
		WithDataset(datasetID).
		WithTasks(tasks).
		WithLabels(projectID).
		Build()
	if err != nil {
		return err
	}
	return s.push(message)
}

func (s *service) pushWithNATS(message PushTaskMessage) error {
	logger.Info("Publishing task...")
	return s.nc.Publish(viper.GetString("annotation.pushtask"), &message)
}

func (s *service) push(message PushTaskMessage) error {
	path := s.annotationBasePath + "/task"
	b, err := json.Marshal(&message)
	if err != nil {
		return err
	}
	logger.Infof("msg: %s", b)
	resp, err := s.client.Post(path, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	bb, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		logger.Info(string(bb))
	}
	return nil
}

type builder struct {
	errs          []error
	workspaceRepo workspace.Repository
	projectRepo   project.Repository
	datasetRepo   dataset.Repository
	labelRepo     label.Repository

	workspace WorkspaceObject
	project   ProjectObject
	dataset   DatasetObject
	tasks     []TaskObject
	labels    []LabelObject
}

func NewBuilder(workspaceRepo workspace.Repository,
	projectRepo project.Repository,
	datasetRepo dataset.Repository,
	labelRepo label.Repository) *builder {
	return &builder{
		errs:          make([]error, 0),
		workspaceRepo: workspaceRepo,
		projectRepo:   projectRepo,
		datasetRepo:   datasetRepo,
		labelRepo:     labelRepo,
	}
}

func (b *builder) WithWorkspace(id uint64) *builder {
	w, err := b.workspaceRepo.Get(id)
	if err != nil {
		logger.Errorf("error getting workspace %d. err %v", id, err)
		b.errs = append(b.errs, err)
		return b
	}
	var admin uint64
	perms, _, err := b.workspaceRepo.GetPermission(w.ID, workspace.Admin, 0, 1)
	if err != nil || len(perms) == 0 {
		b.errs = append(b.errs, errors.WorkspacePermissionErrorCreating.NewWithMessage(""))
	} else {
		admin = perms[0].UserID
	}
	b.workspace = WorkspaceObject{
		ID:    w.ID,
		Title: w.Title,
		Admin: admin,
	}
	return b
}

func (b *builder) WithProject(id uint64) *builder {
	p, err := b.projectRepo.Get(id)
	if err != nil {
		logger.Errorf("error getting project %d. err %v", id, err)
		b.errs = append(b.errs, err)
		return b
	}
	var manager = make([]uint64, 0)
	perms, _, err := b.projectRepo.GetProjectPermissions(p.ID, project.Manager, 0, 1)
	if err != nil {
		logger.Errorf("error getting manager of project %d,. err %v", id, err)
		b.errs = append(b.errs, errors.WorkspacePermissionErrorCreating.NewWithMessage(""))
		return b
	} else {
		for i := range perms {
			manager = append(manager, perms[i].UserID)
		}
	}
	b.project = ProjectObject{
		ID:             p.ID,
		Title:          p.Title,
		ProjectManager: manager,
	}
	return b
}

func (b *builder) WithDataset(id uint64) *builder {
	d, err := b.datasetRepo.Get(id)
	if err != nil {
		logger.Errorf("error get dataset %d. err %v", id, err)
		b.errs = append(b.errs, err)
		return b
	}
	b.dataset = DatasetObject{
		ID:        d.ID,
		Title:     d.Title,
		ProjectID: d.ProjectID,
	}
	return b
}

func (b *builder) WithTasks(tasks []task.Task) *builder {
	var t = make([]TaskObject, len(tasks))
	for i, task := range tasks {
		t[i] = TaskObject{
			ID:        task.ID,
			Labeler:   task.Labeler,
			Reviewer:  task.Reviewer,
			CreatedAt: clock.UnixMillisecondFromTime(task.CreatedAt),
		}
	}
	b.tasks = t
	return b
}

func (b *builder) WithLabels(projectID uint64) *builder {
	labels, err := b.labelRepo.GetByProjectId(projectID)
	if err != nil {
		logger.Errorf("error get project %d. err %v", projectID, err)
		b.errs = append(b.errs, err)
		return b
	}
	var responses = make([]LabelObject, len(labels))
	for i, label := range labels {
		responses[i] = LabelObject{
			ID:    label.ID,
			Name:  label.Name,
			Color: label.Color,
			Tool: ToolObject{
				ID:   label.Tool.ID,
				Name: label.Tool.Name,
			},
		}
	}
	b.labels = responses
	return b
}
func (b *builder) Build() (PushTaskMessage, error) {
	if len(b.errs) != 0 {
		logger.Errorf("error creating %d task %v", len(b.errs), b.errs)
		return PushTaskMessage{}, errors.TaskCannotCreate.NewWithMessage("cannot build push task message")
	}
	return PushTaskMessage{
		Workspace: b.workspace,
		Project:   b.project,
		Dataset:   b.dataset,
		Tasks:     b.tasks,
		Labels:    b.labels,
	}, nil
}

func (s *service) UpdateProject(projectID uint64) error {
	p, err := s.projectRepo.Get(projectID)
	if err != nil {
		logger.Errorf("error getting project %d. err %v", projectID, err)
		return err
	}
	var managers = make([]uint64, 0)
	perms, _, err := s.projectRepo.GetProjectPermissions(p.ID, project.Manager, 0, 1)
	if err != nil {
		logger.Errorf("error getting managers of project %d,. err %v", projectID, err)
		return err
	} else {
		for i := range perms {
			managers = append(managers, perms[i].UserID)
		}
	}
	object := ProjectObject{
		ID:             p.ID,
		Title:          p.Title,
		ProjectManager: managers,
	}
	b, err := json.Marshal(object)
	if err != nil {
		return errors.AnnotationCannotReadBody.NewWithMessage("error marshalling object")
	}
	path := s.annotationBasePath + "/annotation/project/update"
	logger.Infof("publishing project to annotation server. path: %s. body %s", path, b)
	resp, err := s.client.Post(path, "application/json", bytes.NewReader(b))
	if err != nil {
		return errors.AnnotationCannotGetFromServer.NewWithMessageF("error requesting to annotation server. err %v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.AnnotationCannotReadBody.NewWithMessageF("cannot read body from annotation server. err", err)
	}
	logger.Infof("update to annotation server resp %s", body)
	return nil
}

func (s *service) UpdateDataset(datasetID uint64) error {
	d, err := s.datasetRepo.Get(datasetID)
	if err != nil {
		logger.Errorf("[ANNOTATION] - error getting project %d. err %v", datasetID, err)
		return err
	}
	object := DatasetObject{
		ID:        d.ID,
		Title:     d.Title,
		ProjectID: d.ProjectID,
	}
	b, err := json.Marshal(object)
	if err != nil {
		return errors.AnnotationCannotReadBody.NewWithMessage("error marshalling object")
	}
	path := s.annotationBasePath + "/annotation/dataset/update"
	logger.Infof("[ANNOTATION] - publishing project to annotation server. path: %s. body %s", path, b)
	resp, err := s.client.Post(path, "application/json", bytes.NewReader(b))
	if err != nil {
		return errors.AnnotationCannotGetFromServer.NewWithMessageF("error requesting to annotation server. err %v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.AnnotationCannotReadBody.NewWithMessageF("cannot read body from annotation server. err", err)
	}
	logger.Infof("[ANNOTATION] - update to annotation server resp %s", body)
	return nil
}

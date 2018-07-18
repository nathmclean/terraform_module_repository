package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

const (
	ErrNameSpaceRequired modelError = "models: Namespace is required"
	ErrNameRequired      modelError = "models: Name is required"
	ErrVersionRequired   modelError = "models: Version is required"
	ErrProviderRequired  modelError = "models: Provider is required"
	ErrModuleExists      modelError = "models: Module already exists"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

type Module struct {
	gorm.Model  `json:"-"`
	ModuleId    string    `json:"id" gorm:"not_null;UNIQUE;index"`
	Owner       string    `json:"owner"`
	Namespace   string    `gorm:"not_null;index" json:"namespace"`
	Name        string    `gorm:"not_null" json:"name"`
	Version     string    `gorm:"not_null" json:"version"`
	Provider    string    `gorm:"not_null;index" json:"provider"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	Downloads   uint      `json:"downloads"`
	Verified    bool      `json:"verified"`
}

type ModuleDB interface {
	List(request *ListRequest) ([]Module, error)
	ListWithParams(request *ListRequest) ([]Module, error)
	Create(modules *Module) error
	ByModuleId(moduleId string) (*Module, error)
	Get(request *GetRequest)(*Module, error)
}

type ModuleService interface {
	ModuleDB
}

type moduleService struct {
	ModuleDB
}

func NewModuleService(db *gorm.DB) ModuleService {
	return &moduleService{
		&moduleValidator{
			ModuleDB: &moduleGorm{db},
		},
	}
}

type moduleGorm struct {
	db *gorm.DB
}

func (m *moduleGorm) Create(module *Module) error {
	return m.db.Create(module).Error
}

func (m *moduleGorm) Get(request *GetRequest)(*Module, error) {

	return nil, nil
}

// List returns a list of all modules, within the limits of the offset and limit
func (m *moduleGorm) List(request *ListRequest) ([]Module, error) {
	fmt.Println("models List")
	var modules []Module
	db := m.db
	if err := db.Where("id in (?)", db.Table("modules").Select("MAX(id)").Group("name, namespace, provider").Limit(request.Limit).Offset(request.Offset).QueryExpr()).Find(&modules).Error; err != nil {
		fmt.Println("model", err)
		return nil, err
	}
	fmt.Println("modules", modules)
	return modules, nil
}

// ListWithParams providers a request filtered by some or all of namespace, verified, provider and version.
func (m *moduleGorm) ListWithParams(request *ListRequest) ([]Module, error) {
	fmt.Println("models List")
	var modules []Module
	db := m.db
	base := db.Table("modules").Select("MAX(id)")
	var query []string
	var args []interface{}
	if request.Namespace != "" {
		query = append(query, "namespace = ?")
		args = append(args, request.Namespace)
	}
	if request.Name != "" {
		query = append(query, "name = ?")
		args = append(args, request.Name)
	}
	if request.Provider != "" {
		query = append(query, "provider = ?")
		args = append(args, request.Provider)
	}
	where := base
	if len(query) > 0 {
		where = where.Where(strings.Join(query, " AND "), args...)
	}

	grouped := where.Group("name, namespace, provider").Limit(request.Limit).Offset(request.Offset).QueryExpr()
	if err := db.Where("id in (?)", grouped).Find(&modules).Error; err != nil {
		fmt.Println("model", err)
		return nil, err
	}
	fmt.Println("modules", modules)
	return modules, nil
}

func (m *moduleGorm) ByModuleId(moduleId string) (*Module, error) {
	var module Module
	err := m.db.Where("module_id = ?", moduleId).First(&module).Error
	if err != nil {
		return nil, err
	}
	return &module, nil

}

type moduleValFn func(db *Module) error

func runModuleValFns(module *Module, fns ...moduleValFn) error {
	for _, fn := range fns {
		if err := fn(module); err != nil {
			return err
		}
	}
	return nil
}

type moduleValidator struct {
	ModuleDB
}

func (mv *moduleValidator) Create(module *Module) error {
	err := runModuleValFns(module,
		mv.nameRequired,
		mv.namespaceRequired,
		mv.versionRequired,
		mv.providerRequired,
		mv.generateModuleId,
		mv.uniqueModuleID)

	if err != nil {
		return err
	}
	return mv.ModuleDB.Create(module)
}

func (mv *moduleValidator) Get(request *GetRequest)(*Module, error) {

	return nil, nil
}

func (mv *moduleValidator) List(request *ListRequest) ([]Module, error) {
	request.validatePagination()

	return mv.ModuleDB.List(request)
}

func (mv *moduleValidator) ListWithParams(request *ListRequest) ([]Module, error) {
	request.validatePagination()
	return mv.ModuleDB.ListWithParams(request)
}

func (mv *moduleValidator) namespaceRequired(module *Module) error {
	if module.Namespace == "" {
		return ErrNameSpaceRequired
	}
	return nil
}

func (mv *moduleValidator) nameRequired(module *Module) error {
	if module.Name == "" {
		return ErrNameRequired
	}
	return nil
}

func (mv *moduleValidator) versionRequired(module *Module) error {
	if module.Version == "" {
		return ErrVersionRequired
	}
	return nil
}

func (mv *moduleValidator) providerRequired(module *Module) error {
	if module.Provider == "" {
		return ErrProviderRequired
	}
	return nil
}

func (mv *moduleValidator) generateModuleId(module *Module) error {
	// namespace/name/provider/version
	module.ModuleId = fmt.Sprintf("%s/%s/%s/%s", module.Namespace, module.Name, module.Provider, module.Version)
	return nil
}

func (mv *moduleValidator) uniqueModuleID(module *Module) error {
	tmpModule, err := mv.ModuleDB.ByModuleId(module.ModuleId)
	if err != nil {
		switch err.Error() {
		case "record not found":
		default:
			fmt.Println("model uniqueModuleID", err)
			return err
		}
	}
	if tmpModule != nil {
		return ErrModuleExists
	}
	return nil
}

const ListLimit = 20

type ListRequest struct {
	Limit     int
	Offset    int
	Namespace string
	Verified  bool
	Name      string
	Provider  string
	Version   string
}

type GetRequest struct {
	Namespace string
	Name      string
	Provider  string
	Version   string
}

func (p *ListRequest) validatePagination() {
	if p.Offset < 0 {
		p.Offset = 0
	}
	if p.Limit > ListLimit || p.Limit <= 0 {
		p.Limit = ListLimit
	}
}

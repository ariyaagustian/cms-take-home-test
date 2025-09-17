package seed

import (
	"cms/server/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	// baseline roles
	roles := []model.Role{{Name:"Admin"},{Name:"Editor"},{Name:"Viewer"}}
	for _, r := range roles {
		db.FirstOrCreate(&model.Role{}, r)
	}
	// admin user
	pw, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := model.User{Name:"Admin", Email:"admin@cms.local", PasswordHash:string(pw)}
	db.FirstOrCreate(&model.User{}, model.User{Email: admin.Email})
	db.Model(&model.User{}).Where("email = ?", admin.Email).Updates(admin)

	// map admin role
	var u model.User
	db.Where("email = ?", admin.Email).First(&u)
	var role model.Role
	db.Where("name = ?", "Admin").First(&role)
	db.FirstOrCreate(&model.UserRole{}, model.UserRole{UserID: u.ID, RoleID: role.ID})

	// sample content type "post"
	ct := model.ContentType{Name:"Post", Slug:"post"}
	db.FirstOrCreate(&model.ContentType{}, model.ContentType{Slug: ct.Slug})
	db.Model(&model.ContentType{}).Where("slug = ?", ct.Slug).Updates(ct)

	return nil
}

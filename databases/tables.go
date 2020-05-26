package databases

// User 用户表
type User struct {
	ID         int    `gorm:"primary_key;AUTO_INCREMENT;unique_index"` //用户ID 平台唯一ID
	Name       string `gorm:"not null;size:100;unique"`                //登录名
	Password   string `gorm:"not null;size:32"`                        // 用户密码
	Nickname   string `gorm:"not null;size:40"`                        // 用户名称
	Roles      int    `gorm:"not null"`                                // 用户所属角色
	Status     int    `gorm:"not null"`                                // 状态 1启用 2禁用 3 已删除
	Createtime int64  `gorm:"not null"`                                // 创建时间
	Updatetime int64  `gorm:"not null;default:0"`                      //更新时间
	Phone      string `gorm:"not null;default:'';size:11"`             // 电话
	Email      string `gorm:"not null;default:'';size:100"`            // 邮箱
	Address    string `gorm:"not null;default:'';size:255"`            //地址
	Details    string `gorm:"not null;default:'';size:255"`            //备注
	Userid     int    `gorm:"not null;default:0"`                      //创建人员
}

// Grade 等级表
type Grade struct {
	ID         int    `gorm:"primary_key;AUTO_INCREMENT;unique_index"` //权限ID
	Name       string `gorm:"not null;size:100"`                       //角色名称
	Details    string `gorm:"not null;default:'';size:255"`            //备注
	Createtime int64  `gorm:"not null"`                                //创建时间
	Updatetime int64  `gorm:"not null;default:0"`                      //更新时间
	Userid     int    `gorm:"not null;default:0"`                      //创建人员
}

// Roles 角色表
type Roles struct {
	ID         int    `gorm:"primary_key;AUTO_INCREMENT;unique_index"` //权限ID
	Gradeid    int    `gorm:"not null;index"`                          //所属等级
	Name       string `gorm:"not null;size:100"`                       //角色名称
	Details    string `gorm:"not null;default:'';size:255"`            //备注
	Createtime int64  `gorm:"not null"`                                // 创建时间
	Updatetime int64  `gorm:"not null;default:0"`                      //更新时间
	Userid     int    `gorm:"not null;default:0"`                      //创建/更新人员
}

// Application 应用
type Application struct {
	ID         int    `gorm:"primary_key;AUTO_INCREMENT;unique_index"` //应用ID
	Name       string `gorm:"not null;size:100;index"`                 //应用名称
	Details    string `gorm:"not null;default:'';size:255"`            //备注
	Createtime int64  `gorm:"not null"`                                //注册时间
}

// Appauthority 应用权限
type Appauthority struct {
	ID         int    `gorm:"primary_key;AUTO_INCREMENT;unique_index"` //应用权限ID
	Name       string `gorm:"not null;size:100"`                       //应用权限名称
	Details    string `gorm:"not null;default:'';size:255"`            //备注
	Createtime int64  `gorm:"not null"`                                //注册时间
}

// Authority 角色权限
type Authority struct {
	ID             int   `gorm:"primary_key;AUTO_INCREMENT;unique_index"` //角色权限ID
	Rolesid        int   `gorm:"not null;index"`                          //角色ID
	Applicationid  int   `gorm:"not null;index"`                          //应用ID
	Appauthorityid int   `gorm:"not null;index"`                          //应用权限ID
	Authority      int   `gorm:"not null"`                                //授权状况 2授权 其他未授权
	Createtime     int64 `gorm:"not null"`                                // 创建时间
	Updatetime     int64 `gorm:"not null;default:0"`                      //更新时间
	Userid         int   `gorm:"not null;default:0"`                      //创建人员
}

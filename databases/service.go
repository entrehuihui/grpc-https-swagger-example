package databases

// Create 创建
func Create(table interface{}) error {
	return db.Create(table).Error
}

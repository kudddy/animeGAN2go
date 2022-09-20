package pg

func InsertCancelAction(userid int) {
	result := &UserFileIdStatus{}
	db.Table("user_file_id_status").Where("user_id = ?", userid).Delete(result)
}

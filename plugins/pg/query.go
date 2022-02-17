package pg

func InsertCancelAction(userid int) {
	result := &UserFileIdStatus{}
	//GetDB().Table("job_canceled_actions").Update(result)
	//GetDB().Table("user_file_id_status").Where("user_id = ?", userid).Delete(result)
	postgres.Table("user_file_id_status").Where("user_id = ?", userid).Delete(result)

}

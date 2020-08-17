package models

type Visitor struct {
	Model
	Name string `json:"name"`
	Avatar string `json:"avatar"`
	SourceIp string `json:"source_ip"`
	ToId string `json:"to_id"`
	VisitorId string `json:"visitor_id"`
	Status uint `json:"status"`
	Refer string `json:"refer"`
	City string `json:"city"`
	ClientIp string `json:"client_ip"`
}
func CreateVisitor(name string,avatar string,sourceIp string,toId string,visitorId string,refer string,city string,clientIp string){
	old:=FindVisitorByVisitorId(visitorId)
	if old.Name!=""{
		//更新状态上线
		UpdateVisitorStatus(visitorId,1)
		return
	}
	v:=&Visitor{
		Name:name,
		Avatar: avatar,
		SourceIp:sourceIp,
		ToId:toId,
		VisitorId: visitorId,
		Status:1,
		Refer:refer,
		City:city,
		ClientIp:clientIp,
	}
	DB.Create(v)
}
func FindVisitorByVisitorId(visitorId string)Visitor{
	var v Visitor
	DB.Where("visitor_id = ?", visitorId).First(&v)
	return v
}
func FindVisitors(page uint,pagesize uint)[]Visitor{
	offset:=(page-1)*pagesize
	if offset<0{
		offset=0
	}
	var visitors []Visitor
	DB.Offset(offset).Limit(pagesize).Order("status desc, updated_at desc").Find(&visitors)
	return visitors
}
func FindVisitorsByKefuId(page uint,pagesize uint,kefuId string)[]Visitor{
	offset:=(page-1)*pagesize
	if offset<0{
		offset=0
	}
	var visitors []Visitor
	DB.Where("to_id=?",kefuId).Offset(offset).Limit(pagesize).Order("status desc, updated_at desc").Find(&visitors)
	return visitors
}
func FindVisitorsOnline()[]Visitor{
	var visitors []Visitor
	DB.Where("status = ?",1).Find(&visitors)
	return visitors
}
func UpdateVisitorStatus(visitorId string,status uint){
	visitor:=Visitor{
	}
	DB.Model(&visitor).Where("visitor_id = ?",visitorId).Update("status", status)
}
//查询条数
func CountVisitors()uint{
	var count uint
	DB.Model(&Visitor{}).Count(&count)
	return count
}
//查询条数
func CountVisitorsByKefuId(kefuId string)uint{
	var count uint
	DB.Model(&Visitor{}).Where("to_id=?",kefuId).Count(&count)
	return count
}
package main

import (
	"fmt"
	"time"

	"deadlock/common/govalidator"
	fn "tericai/common/helpers"
	md "tericai/common/models"

	"github.com/globalsign/mgo/bson"
)

func projectList(req fn.Request) {
	defer fn.RecoverPanic()

	rules := govalidator.MapData{
		"page":  []string{"numeric_between:1,999999"},
		"limit": []string{"numeric_between:1,999999"},
	}

	isValid := fn.ValidateRespond(req, rules, collections)

	if isValid {
		var results fn.PaginateResult
		var filter = make(bson.M)

		//// filter results by name or email or phone
		//if req.Body["role"] != nil {
		//	filter["role"] = req.Body["role"].(string)
		//}

		err := fn.Paginate(collections["Project"], filter, nil, req.Body["page"], req.Body["limit"], "", &results)

		if err != nil {
			resp := fn.Response{Error: "pagination error", Exception: err, HTTPCode: 500}
			fn.ReplyBack(req, resp)
			return
		}

		if err != nil {
			resp := fn.Response{Error: "cannot get data", Exception: err, Data: nil, HTTPCode: 500}
			fn.ReplyBack(req, resp)
		} else {
			resp := fn.Response{Msg: "projects list", Data: results, HTTPCode: 200}
			fn.ReplyBack(req, resp)
		}
	}

}

func projectAdd(req fn.Request) {

	defer fn.RecoverPanic()

	rules := govalidator.MapData{
		"title":       []string{"required"},
		"description": []string{},
	}

	isValid := fn.ValidateRespond(req, rules, collections)

	if isValid {

		c := md.ProjectModel{
			ID:          bson.NewObjectId(),
			Title:       req.Body["title"].(string),
			Description: req.Body["description"].(string),
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		fmt.Println(" i am here")

		err := collections["Project"].Insert(c)
		if err != nil {
			fmt.Println(" Log -- err", err)
			resp := fn.Response{Error: "cannot insert project", Exception: err, HTTPCode: 500}
			fn.ReplyBack(req, resp)
			return
		}

		resp := fn.Response{Msg: "project inserted", Data: c, HTTPCode: 201}
		fn.ReplyBack(req, resp)
	}
}

func projectEdit(req fn.Request) {
	defer fn.RecoverPanic()

	rules := govalidator.MapData{
		"project_id": []string{"required", "objectID", "exists:Project,_id"},
		"is_active":  []string{"bool"},
	}

	isValid := fn.ValidateRespond(req, rules, collections)

	if isValid {

		projectID := req.Body["project_id"].(string)
		var project md.ProjectModel

		collections["Project"].FindId(bson.ObjectIdHex(projectID)).One(&project)

		if req.Body["title"] != nil {
			project.Title = req.Body["title"].(string)
		}

		if req.Body["description"] != nil {
			project.Description = req.Body["description"].(string)
		}

		project.UpdatedAt = time.Now()

		err2 := collections["Project"].Update(bson.M{"_id": project.ID}, bson.M{"$set": project})
		if err2 != nil {
			resp := fn.Response{Error: "cannot edit project", Exception: err2, HTTPCode: 500}
			fn.ReplyBack(req, resp)
			return
		}

		resp := fn.Response{Msg: "edited project", HTTPCode: 201}
		fn.ReplyBack(req, resp)

	}
}

func projectDelete(req fn.Request) {
	defer fn.RecoverPanic()

	rules := govalidator.MapData{
		"project_id": []string{"required", "objectID", "exists:Project,_id"},
	}

	isValid := fn.ValidateRespond(req, rules, collections)

	if isValid {
		projectID := req.Body["project_id"].(string)
		err := collections["Project"].Remove(bson.M{"_id": bson.ObjectIdHex(projectID)})

		if err != nil {
			resp := fn.Response{Error: "cannot delete", Exception: err, HTTPCode: 500}
			fn.ReplyBack(req, resp)
			return
		}

		resp := fn.Response{Msg: "project deleted successfully", HTTPCode: 200}
		fn.ReplyBack(req, resp)
	}
}

func projectView(req fn.Request) {
	defer fn.RecoverPanic()

	rules := govalidator.MapData{
		"project_id": []string{"required", "objectID"},
	}

	isValid := fn.ValidateRespond(req, rules, collections)

	if isValid {
		projectID := req.Body["project_id"].(string)
		var project md.ProjectModel
		err := collections["Project"].FindId(bson.ObjectIdHex(projectID)).One(&project)

		if err != nil {
			resp := fn.Response{Error: "cannot view project data", Exception: err, HTTPCode: 500}
			fn.ReplyBack(req, resp)
			return
		}

		resp := fn.Response{Msg: "project data", Data: project, HTTPCode: 200}
		fn.ReplyBack(req, resp)
	}

}

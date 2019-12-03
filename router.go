package main

import (
	"tericai/common/helpers"
	"fmt"
)

func callMethod(callbackFn string, req helpers.Request) {
	fmt.Println("calling -->", callbackFn)

	if callbackFn == "project_add" {
		projectAdd(req)
	} else if callbackFn == "project_list" {
		projectList(req)
	} else if callbackFn == "project_edit" {
		projectEdit(req)
	} else if callbackFn == "project_delete" {
		projectDelete(req)
	} else if callbackFn == "project_view" {
		projectView(req)
	}

}

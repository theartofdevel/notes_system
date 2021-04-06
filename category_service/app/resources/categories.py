import json
import logging
from http import HTTPStatus

from flask import make_response, jsonify, request
from flask_apispec import marshal_with, use_kwargs, MethodResource
from flask_restful import Resource
from injector import inject

from dao.model.dto import CreateCategoryDTO, UpdateCategoryDTO, DeleteCategoryDTO
from dao.model.model import Category
from helpers.json import CustomJSONEncoder
from service import CategoryService


class CategoryResource(MethodResource, Resource):
    __slots__ = ["service", "logger"]

    @inject
    def __init__(self, service: CategoryService, logger: logging.Logger):
        self.service = service
        self.logger = logger

    @use_kwargs(UpdateCategoryDTO.to_schema())
    @marshal_with(None, code=HTTPStatus.NO_CONTENT)
    def patch(self, category_dto: UpdateCategoryDTO, cuuid: str):
        category_dto.uuid = cuuid
        self.service.update_category(category=category_dto)
        r = make_response(jsonify(), HTTPStatus.NO_CONTENT)
        r.headers["Content-Type"] = "application/json"
        return r

    @use_kwargs(DeleteCategoryDTO.to_schema())
    @marshal_with(None, code=HTTPStatus.NO_CONTENT)
    def delete(self, category_dto: DeleteCategoryDTO, cuuid: str):
        category_dto.uuid = cuuid
        self.service.delete_category(category=category_dto)
        r = make_response(jsonify(), HTTPStatus.NO_CONTENT)
        r.headers["Content-Type"] = "application/json"
        return r


class CategoriesResource(MethodResource, Resource):
    __slots__ = ["service", "logger"]

    @inject
    def __init__(self, service: CategoryService, logger: logging.Logger):
        self.service = service
        self.logger = logger

    @marshal_with(Category.to_schema())
    def get(self):
        user_uuid = request.args.get("user_uuid")
        categories = self.service.get_categories(user_uuid=user_uuid)
        r = make_response(json.dumps(categories, cls=CustomJSONEncoder), HTTPStatus.OK)
        r.headers["Content-Type"] = "application/json"
        return r

    @use_kwargs(CreateCategoryDTO.to_schema())
    @marshal_with(None, code=HTTPStatus.NO_CONTENT)
    def post(self, category_dto: CreateCategoryDTO):
        category = self.service.create_category(category=category_dto)
        r = make_response(jsonify(), HTTPStatus.NO_CONTENT)
        r.headers["Location"] = f"/api/categories/{category.uuid}"
        r.headers["Content-Type"] = "application/json"
        return r

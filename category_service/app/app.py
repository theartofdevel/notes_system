import logging
import os
from pprint import pprint

from flask import Flask, make_response, jsonify
from flask_cors import CORS
from flask_injector import FlaskInjector
from flask_restful import Api
from injector import Injector
from webargs.flaskparser import parser, abort

from config import Config
from constants import LOG_DIR, CONFIG_FILE_PATH, ANY_ORIGIN, EXPOSE_HEADERS
from di import StorageModule, LoggerModule
from exceptions import AppException, ValidationException, AppError
from helpers.flask import app_exception_handler, uncaught_exception_handler
from resources import CategoryResource, CategoriesResource

config = Config(yaml_file=CONFIG_FILE_PATH)

logger = logging.getLogger("main")
logger.setLevel(logging.DEBUG)

if not os.path.exists(LOG_DIR):
    os.makedirs(LOG_DIR)

fh = logging.FileHandler(f"{LOG_DIR}/all.log")
fh.setLevel(logging.DEBUG)
formatter = logging.Formatter("%(asctime)s - %(name)s - %(levelname)s - %(message)s")
fh.setFormatter(formatter)
logger.addHandler(fh)

app = Flask(__name__)
app.url_map.strict_slashes = False
app.config.from_object(config)

api = Api(app)
api.add_resource(CategoriesResource, "/api/categories")
api.add_resource(CategoryResource, "/api/categories/<string:cuuid>")

CORS(app, resources={"*": ANY_ORIGIN}, expose_headers=EXPOSE_HEADERS)

injector = Injector([StorageModule(config=config), LoggerModule(logger=logger)])
FlaskInjector(app=app, injector=injector)

app.errorhandler(AppException)(app_exception_handler)

if not config.DEBUG:
    app.errorhandler(Exception)(uncaught_exception_handler)


@parser.error_handler
def handle_request_parsing_error(err, req, schema, *, error_status_code, error_headers):
    developer_message = ""
    for field_name, errors in err.messages["json"].items():
        developer_message += f"Error field '{field_name}': "
        c = 0
        for i in errors:
            c += 1
            developer_message += f"{i}"
            if c != len(errors):
                developer_message += ","
    raise ValidationException(exc_data=AppError.VALIDATION_ERROR, developer_message=developer_message)


pprint(app.url_map)

if __name__ == '__main__':
    app.run(host="localhost", port=10001, debug=True)

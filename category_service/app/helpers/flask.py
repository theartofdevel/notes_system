import logging
from traceback import format_exc

from flask import make_response, jsonify

from exceptions import AppException, AppError


def app_exception_handler(exception):
    r = make_response(
        jsonify(exception.__dict__), exception.http_code
    )
    return r


def uncaught_exception_handler(exception):
    logging.getLogger("main").error(f"* Uncaught exception [{exception}]: {format_exc}")
    return app_exception_handler(AppException(exc_data=AppError.SYSTEM_ERROR))

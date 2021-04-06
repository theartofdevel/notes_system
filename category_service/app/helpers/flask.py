import logging
import sys
import traceback
from http import HTTPStatus
from traceback import format_exc

from flask import make_response, jsonify

from exceptions import AppException, AppError, NotFoundException, ValidationException


def app_exception_handler(exception):
    http_code = 418
    if isinstance(exception, NotFoundException):
        http_code = HTTPStatus.NOT_FOUND
    elif isinstance(exception, ValidationException):
        http_code = HTTPStatus.BAD_REQUEST
    r = make_response(
        jsonify(exception.__dict__), http_code
    )
    return r


def uncaught_exception_handler(exception: Exception):
    # Get current system exception
    ex_type, ex_value, ex_traceback = sys.exc_info()

    # Extract unformatter stack traces as tuples
    trace_back = traceback.extract_tb(ex_traceback)

    # Format stacktrace
    stack_trace = list()

    for trace in trace_back:
        stack_trace.append("File : %s , Line : %d, Func.Name : %s, Message : %s" % (trace[0], trace[1], trace[2], trace[3]))

    logging.getLogger("main").error(f"* Uncaught exception [{exception}]: {format_exc}")
    return app_exception_handler(AppException(exc_data=AppError.SYSTEM_ERROR, developer_message=str(
        exception)))

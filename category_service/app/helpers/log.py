import functools
import logging
import sys


def logged(func):
    @functools.wraps(func)
    def sync_wrapper(*args, **kwargs):
        caller = sys._getframe(1).f_code.co_name
        caller_filename = sys._getframe(1).f_code.co_filename

        cfs = caller_filename.split("/")
        if len(cfs) > 4:
            caller_filename = "/".join(cfs[-3:])

        args_str = None
        if len(args) > 1:
            args_str = ", ".join([str(arg) for arg in args[1:]])
        kwargs_str = None
        if len(kwargs) > 0:
            kwargs_str = ", ".join([":".join([str(j) for j in i]) for i in kwargs.items()])

        result = func(*args, **kwargs)

        clazz_func_name = "{}.{}.{}".format(func.__module__, args[0].__class__.__name__, func.__name__)
        log_str = f"{caller_filename}:{caller}() -> {clazz_func_name}()"
        if args_str:
            log_str += f" args: {args_str}"

        if kwargs_str:
            log_str += f" kwargs: {kwargs_str}"
        if result:
            log_str += f" result: {result}"

        logging.getLogger("main").debug(log_str)
        return result

    return sync_wrapper
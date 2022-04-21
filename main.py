#!/usr/bin/env python3

from contextlib import contextmanager

import time


class FibGenerator():
    def __init__(self, limit=100):
        self.limit = limit
        self.n = 0
        self.m = 1

    def get_next(self):
        res = self.n
        self.n = self.m
        self.m = self.m + res
        return res


def fib(limit: int=100):
    n, m = 0, 1
    for _ in range(limit):
        yield n
        n, m = m, n+m
    return


@contextmanager
def timer(task: str):
    start = time.time()
    try:
        yield
    except Exception as ex:
        print(f'time elapsed before failing {task}: {time.time() - start}')
        raise ex
    print(f'time elapsed for {task}: {time.time() - start}')

@timer("print_fib")
def print_fib():
    fib_gen = fib()
    fib_gen.limit = 1000
    for n in fib_gen:
        print(n)


@timer("main")
def main():
    print_fib()


if __name__ == '__main__':
    main()

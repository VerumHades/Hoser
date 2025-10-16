import os
from sys import executable
from subprocess import Popen, CREATE_NEW_CONSOLE

def spawn_function_shell(function):
    import_template = f"from dev import {function.__name__} \n{function.__name__}()"
    Popen([executable, '-c', import_template], creationflags=CREATE_NEW_CONSOLE)


def start_go_server():
    print('A')
    os.chdir("server")
    os.system("air")
    input()
    
def start_vite():
    os.chdir("client")
    os.system("npm run dev")
    input()

if __name__ == "__main__":
    spawn_function_shell(start_go_server)
    spawn_function_shell(start_vite)
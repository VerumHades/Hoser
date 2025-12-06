import os
import sys
from subprocess import Popen

def spawn_function_shell(function):
    """
    Spawns a new terminal window and runs the given function from the 'dev' module.
    Works on both Windows and Linux.
    """
    function_name = function.__name__
    import_template = f"from dev import {function_name}; {function_name}()"

    if os.name == "nt":  # Windows
        from subprocess import CREATE_NEW_CONSOLE
        # Use CREATE_NEW_CONSOLE to open a new terminal window
        Popen([sys.executable, "-c", import_template], creationflags=CREATE_NEW_CONSOLE)

    elif sys.platform.startswith("linux"):  # Linux
        # Use gnome-terminal to spawn a new terminal window
        command = f"{sys.executable} -c '{import_template}'"
        Popen(["gnome-terminal", "--", "bash", "-c", command])

    else:
        raise OSError("Unsupported operating system for terminal spawning")

# Example usage
def start_go_server():
    print("A")
    os.chdir(os.path.join("services","api"))
    os.system("air")
    input()

def start_vite():
    os.chdir(os.path.join("services", "frontend", "client"))
    os.system("npm run dev")
    input()

if __name__ == "__main__":
    spawn_function_shell(start_go_server)
    spawn_function_shell(start_vite)

package game

import (
	"fmt"
	"log"
	"strings"

	"linuxmodule/pkg/container"
)

type Task struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Hint           string `json:"hint"`
	ValidationType string `json:"validation_type"` // e.g., "file_exists", "dir_exists", "permission_check", "file_contains", "input_match", "process_running", "process_not_running", "user_exists", "group_exists", "owner_check", "group_check", "hostname_check"
	Target         string `json:"target"`          // Path to file/folder or identifier
	Expected       string `json:"expected"`        // Expected value (for input_match or file_contains)
}

type Chapter struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Tasks       []Task `json:"tasks"`
}

var Chapters = []Chapter{
	{
		ID:          1,
		Title:       "Chapter 1: Escape the Potions Dungeon (Basics)",
		Description: "You are a student at Codewarts School of Witchcraft and Wizardry. Snape has locked you inside the dungeons after class. Your wand-terminal is your only tool to escape. Use basic command-spells to inspect the room.",
		Tasks: []Task{
			{
				ID:             1,
				Title:          "Who Am I?",
				Description:    "Verify your magical identity on the server. Run `whoami` in the terminal to inspect the active wizard username.",
				Hint:           "Type 'whoami' into the shell and hit enter. The terminal should print your active student name.",
				ValidationType: "input_match",
				Expected:       "player",
			},
			{
				ID:             2,
				Title:          "Locating Your Coordinates",
				Description:    "Identify the absolute path of the dungeon floor you are currently standing on. Run `pwd` to print it.",
				Hint:           "Running 'pwd' (Print Working Directory) outputs the exact file path of the room your shell is looking at.",
				ValidationType: "input_match",
				Target:         "pwd",
				Expected:       "/home/player",
			},
			{
				ID:             3,
				Title:          "Scanning for Hidden Relics",
				Description:    "Dark magic files are hidden as dotfiles. List all files in the dungeon—including hidden ones (starting with a dot '.')—to search for spell scrolls. Run `ls -la` in the terminal.",
				Hint:           "'ls' lists files, '-a' includes hidden files, and '-l' displays long details. Look for a hidden file starting with '.cyber'.",
				ValidationType: "input_match",
				Target:         "hidden_file",
				Expected:       ".cyber_key",
			},
			{
				ID:             4,
				Title:          "Reading the Spell Scroll",
				Description:    "Read the contents of the hidden scroll `.cyber_key` using the `cat` spell command to reveal the dungeon lock password.",
				Hint:           "Use 'cat .cyber_key' to print the file contents in the terminal.",
				ValidationType: "input_match",
				Target:         "code",
				Expected:       "GRID_PASS_99",
			},
			{
				ID:             5,
				Title:          "Creating Your First Scroll",
				Description:    "Write the spell 'Alohomora' into a new file named `scroll.txt` in your home directory.",
				Hint:           "Use redirection: `echo Alohomora > scroll.txt` or edit with nano.",
				ValidationType: "file_contains",
				Target:         "/home/player/scroll.txt",
				Expected:       "Alohomora",
			},
			{
				ID:             6,
				Title:          "Revealing Scroll Details",
				Description:    "Inspect the details of the created scroll. Use `file scroll.txt` to find out what type of file it is, and enter the file type (e.g. 'ASCII text').",
				Hint:           "Run `file scroll.txt` and look at the description. Enter 'ascii text' (case-insensitive).",
				ValidationType: "input_match",
				Expected:       "ascii text",
			},
		},
	},
	{
		ID:          2,
		Title:       "Chapter 2: The Filesystem Catacombs (Navigation & Directories)",
		Description: "You have slipped out of the dungeon and entered the catacombs of the library. To craft your tracking charms, navigate the filesystem and build a secure research folder.",
		Tasks: []Task{
			{
				ID:             1,
				Title:          "Restricted Archives Discovery",
				Description:    "List the folders inside `/var/mainframe` (`ls /var/mainframe`). Identify and enter the name of the folder related to backups or archives.",
				Hint:           "Run 'ls /var/mainframe'. There are three folders listed; choose the one named archives.",
				ValidationType: "input_match",
				Target:         "archive_folder",
				Expected:       "archives",
			},
			{
				ID:             2,
				Title:          "Creating Your Spellbook Directory",
				Description:    "Create a new directory named `grid` inside your home directory (`/home/player`) to hold your custom charms.",
				Hint:           "Use the make directory command: 'mkdir /home/player/grid' or 'mkdir grid' if you are inside your home folder.",
				ValidationType: "dir_exists",
				Target:         "/home/player/grid",
			},
			{
				ID:             3,
				Title:          "Scribing the Exploit Charm",
				Description:    "Create an empty spell script named `hack.sh` inside your new `/home/player/grid` directory.",
				Hint:           "Use the 'touch' command: 'touch /home/player/grid/hack.sh'.",
				ValidationType: "file_exists",
				Target:         "/home/player/grid/hack.sh",
			},
			{
				ID:             4,
				Title:          "Stealing Library Feeds",
				Description:    "Copy the library mainframe log `/var/log/mainframe.log` to your workspace as `system.log` to monitor prefect movements.",
				Hint:           "Use the copy command 'cp': 'cp /var/log/mainframe.log /home/player/grid/system.log'.",
				ValidationType: "file_exists",
				Target:         "/home/player/grid/system.log",
			},
			{
				ID:             5,
				Title:          "Teleporting Scrolls",
				Description:    "Move the `scroll.txt` file from `/home/player` into your new `/home/player/grid` folder using the move command.",
				Hint:           "Use 'mv /home/player/scroll.txt /home/player/grid/scroll.txt'.",
				ValidationType: "file_exists",
				Target:         "/home/player/grid/scroll.txt",
			},
			{
				ID:             6,
				Title:          "Creating a Secret Backup Directory",
				Description:    "Create a nested backup directory named `/home/player/grid/backup/secrets` (using `mkdir -p` or nested `mkdir`).",
				Hint:           "Use: `mkdir -p /home/player/grid/backup/secrets` to create nested folders.",
				ValidationType: "dir_exists",
				Target:         "/home/player/grid/backup/secrets",
			},
		},
	},
	{
		ID:          3,
		Title:       "Chapter 3: The Process Catacombs (Process Beasts)",
		Description: "Year 3 requires you to master process control. Summon background spirits, inspect active daemons, and banish them when they run rogue.",
		Tasks: []Task{
			{
				ID:             1,
				Title:          "Summon Background Spirit",
				Description:    "Run a process in the background. Execute `sleep 999 &` to launch a sleeping spirit process.",
				Hint:           "Type 'sleep 999 &' in the terminal and press Enter. The '&' symbol spawns the process in the background.",
				ValidationType: "process_running",
				Target:         "sleep",
			},
			{
				ID:             2,
				Title:          "Locate Spirit Identity",
				Description:    "List your background jobs or running processes. Enter the name of the executable you just placed in the background.",
				Hint:           "Run 'jobs' or 'ps' to see your background tasks. Enter 'sleep' to confirm.",
				ValidationType: "input_match",
				Expected:       "sleep",
			},
			{
				ID:             3,
				Title:          "Summon a Second Spirit",
				Description:    "Summon another sleeping process running at `sleep 888 &` so that multiple jobs are active in the background.",
				Hint:           "Type 'sleep 888 &' and hit enter. Run 'jobs' to see both tasks running.",
				ValidationType: "process_running",
				Target:         "sleep",
			},
			{
				ID:             4,
				Title:          "Real-time System Scan",
				Description:    "What command is used to display a dynamically updated, real-time list of all active processes and resource usage in Linux?",
				Hint:           "This tool is commonly used to monitor system performance in real-time. Type 'top'.",
				ValidationType: "input_match",
				Expected:       "top",
			},
			{
				ID:             5,
				Title:          "Banish the Process",
				Description:    "Banish all background sleep spirits. Execute the command to terminate all processes named 'sleep' simultaneously.",
				Hint:           "Use the 'killall' command: 'killall sleep'. Verify via 'jobs' that none are left.",
				ValidationType: "process_not_running",
				Target:         "sleep",
			},
			{
				ID:             6,
				Title:          "Finding Spirit Process ID",
				Description:    "What command is used to query the Process ID (PID) of running processes by name directly? (Hint: starts with 'pg').",
				Hint:           "Type 'pgrep'. This command returns the PIDs of processes matching a name pattern.",
				ValidationType: "input_match",
				Expected:       "pgrep",
			},
		},
	},
	{
		ID:          4,
		Title:       "Chapter 4: Blinding the Prefect Wards (Permissions & Redirections)",
		Description: "Prefect wards are locking down the hallway. Modify scroll access permissions, redirect standard streams, and execute your script to blind the magic sensors.",
		Tasks: []Task{
			{
				ID:             1,
				Title:          "Inspect Scroll Permissions",
				Description:    "Inspect the permissions of `/home/player/grid/hack.sh` using `ls -l`. Type the default owner permissions block (e.g. `rw-` or `-rw-r--r--`).",
				Hint:           "Run 'ls -l /home/player/grid/hack.sh'. Look at the first set of characters after the dash: 'rw-'.",
				ValidationType: "input_match",
				Target:         "permissions",
				Expected:       "rw-",
			},
			{
				ID:             2,
				Title:          "Infusing Magic Execution",
				Description:    "Make the spell script `/home/player/grid/hack.sh` executable for the user.",
				Hint:           "Use the change mode command: 'chmod +x /home/player/grid/hack.sh' or 'chmod 755 /home/player/grid/hack.sh'.",
				ValidationType: "permission_check",
				Target:         "/home/player/grid/hack.sh",
			},
			{
				ID:             3,
				Title:          "Casting the Blindness Spell",
				Description:    "Write a command inside `hack.sh` that writes 'BYPASS' to `/home/player/grid/token.txt`. Execute the script with `./hack.sh` inside your grid folder.",
				Hint:           "You can edit the file with nano/vim or append using redirect: 'echo BYPASS > /home/player/grid/token.txt' inside the script. Execute with './hack.sh' or 'bash hack.sh'.",
				ValidationType: "file_contains",
				Target:         "/home/player/grid/token.txt",
				Expected:       "BYPASS",
			},
			{
				ID:             4,
				Title:          "Appending Spell History",
				Description:    "Append the word 'Lumos' to the end of `/home/player/grid/token.txt` using the double redirection symbol (`>>`).",
				Hint:           "Run: `echo Lumos >> /home/player/grid/token.txt` in the terminal.",
				ValidationType: "file_contains",
				Target:         "/home/player/grid/token.txt",
				Expected:       "Lumos",
			},
			{
				ID:             5,
				Title:          "Blending Error Streams",
				Description:    "Attempt to list a non-existent directory `ls non_existent` and redirect only the error stream (stderr) to `/home/player/grid/errors.txt`.",
				Hint:           "Use '2>' to redirect stderr. Run: `ls non_existent 2> /home/player/grid/errors.txt`.",
				ValidationType: "file_exists",
				Target:         "/home/player/grid/errors.txt",
			},
			{
				ID:             6,
				Title:          "Scribing with Standard Input",
				Description:    "Redirect the contents of `/home/player/grid/errors.txt` into a new file named `/home/player/grid/stderr_backup.txt` using standard input redirection '<' combined with 'cat' and '>': `cat < /home/player/grid/errors.txt > /home/player/grid/stderr_backup.txt`.",
				Hint:           "Run: `cat < /home/player/grid/errors.txt > /home/player/grid/stderr_backup.txt`.",
				ValidationType: "file_exists",
				Target:         "/home/player/grid/stderr_backup.txt",
			},
		},
	},
	{
		ID:          5,
		Title:       "Chapter 5: The Network Realms (Networking & Portals)",
		Description: "Travel between school server sectors. Master domain pings, fetch remote assets, and inspect interface connections.",
		Tasks: []Task{
			{
				ID:             1,
				Title:          "Identify Loopback Gateway",
				Description:    "What is the standard local loopback IP address representing 'localhost' in networking?",
				Hint:           "The loopback interface is a virtual network interface that maps back to the local host. Type '127.0.0.1'.",
				ValidationType: "input_match",
				Expected:       "127.0.0.1",
			},
			{
				ID:             2,
				Title:          "Pinging the Local Gateway",
				Description:    "Use the `ping` utility to send exactly 3 packets to localhost: `ping -c 3 127.0.0.1`. What is the count of packets sent?",
				Hint:           "Look at the '-c' argument in your ping command. You specified sending 3 packets.",
				ValidationType: "input_match",
				Expected:       "3",
			},
			{
				ID:             3,
				Title:          "Identify Web Fetcher tool",
				Description:    "Which command line utility is used to transfer data from or to a server using protocols like HTTP, HTTPS, or FTP?",
				Hint:           "This tool is installed inside your container and starts with the letter 'c'. Type 'curl'.",
				ValidationType: "input_match",
				Expected:       "curl",
			},
			{
				ID:             4,
				Title:          "Fetching Server Coordinates",
				Description:    "Use `curl` to download headers from Google and save the output to `/home/player/grid/web.txt`.",
				Hint:           "Run: `curl -I https://www.google.com > /home/player/grid/web.txt`.",
				ValidationType: "file_exists",
				Target:         "/home/player/grid/web.txt",
			},
			{
				ID:             5,
				Title:          "Loopback Interface Name",
				Description:    "What is the short name traditionally assigned to the loopback interface on Linux systems (as shown in `ip link` or `ifconfig`)?",
				Hint:           "This name consists of two letters. Type 'lo'.",
				ValidationType: "input_match",
				Expected:       "lo",
			},
			{
				ID:             6,
				Title:          "Inspecting Magical Host Name",
				Description:    "Cast the `hostname` spell to print the host identifier of this terminal container, and type the exact output here.",
				Hint:           "Type `hostname` in the terminal, copy the printed host identifier, and submit it.",
				ValidationType: "hostname_check",
			},
		},
	},
	{
		ID:          6,
		Title:       "Chapter 6: The Ritual Scripts (Bash Scripting)",
		Description: "Write automation scripts using variables, arguments, conditionals, and loops to cast complex chain rituals.",
		Tasks: []Task{
			{
				ID:             1,
				Title:          "Scribing the Ritual Variable",
				Description:    "Create a file named `variables.sh` inside `/home/player/grid`. Add `SPELL=\"Lumos\"` and run it. Verify that the file exists and contains `SPELL`.",
				Hint:           "Use 'touch /home/player/grid/variables.sh' and write 'SPELL=\"Lumos\"' inside it.",
				ValidationType: "file_contains",
				Target:         "/home/player/grid/variables.sh",
				Expected:       "SPELL",
			},
			{
				ID:             2,
				Title:          "Creating a Loop Script",
				Description:    "Create a file named `loop.sh` inside `/home/player/grid`. Add a loop check (contains the word `for`) and write it. Verify that the file exists and contains `for`.",
				Hint:           "Create the file and write a simple bash loop. Verify the file contains 'for'.",
				ValidationType: "file_contains",
				Target:         "/home/player/grid/loop.sh",
				Expected:       "for",
			},
			{
				ID:             3,
				Title:          "Creating a Conditional Script",
				Description:    "Create a file named `conditional.sh` inside `/home/player/grid`. Write a conditional statement checking a file (contains the word `if`). Verify the file contains `if`.",
				Hint:           "Create the file and add an 'if' check. Verify the file contains 'if'.",
				ValidationType: "file_contains",
				Target:         "/home/player/grid/conditional.sh",
				Expected:       "if",
			},
			{
				ID:             4,
				Title:          "Script Argument Passing",
				Description:    "Create a file named `greet.sh` inside `/home/player/grid` that references the first command-line argument using the variable `$1`.",
				Hint:           "Create the file greet.sh and make sure it has the characters '$1' inside it (e.g. `echo Hello $1`).",
				ValidationType: "file_contains",
				Target:         "/home/player/grid/greet.sh",
				Expected:       "$1",
			},
			{
				ID:             5,
				Title:          "Inspect Exit Codes",
				Description:    "In Bash scripting, what special variable holds the exit status code of the last executed command?",
				Hint:           "This consists of a dollar sign and a question mark. Type '$?'.",
				ValidationType: "input_match",
				Expected:       "$?",
			},
			{
				ID:             6,
				Title:          "Interpreter Directive (Shebang)",
				Description:    "What is the two-character sequence starting with '#' and '!' used at the very beginning of script files to specify the interpreter?",
				Hint:           "This sequence is commonly called a shebang. Type '#!'.",
				ValidationType: "input_match",
				Expected:       "#!",
			},
		},
	},
	{
		ID:          7,
		Title:       "Chapter 7: The Disk & Storage Catacombs (Storage & Links)",
		Description: "Check system filesystem volumes, inspect folder sizes, and create symbolic file portals.",
		Tasks: []Task{
			{
				ID:             1,
				Title:          "Total Volume Scan",
				Description:    "What command displays the total disk space usage of all mounted filesystems in human-readable format?",
				Hint:           "This command stands for 'disk free' and usually has a '-h' flag. Type 'df -h' or 'df'.",
				ValidationType: "input_match",
				Expected:       "df",
			},
			{
				ID:             2,
				Title:          "Inspect Folder Size",
				Description:    "What command estimates file and directory space usage (commonly used with `-sh` for a summary)?",
				Hint:           "This stands for 'disk usage'. Type 'du'.",
				ValidationType: "input_match",
				Expected:       "du",
			},
			{
				ID:             3,
				Title:          "Create Symbolic Link",
				Description:    "Create a symbolic link (shortcut) named `/home/player/grid/reactor_link` that points to `/var/log/reactor.log`.",
				Hint:           "Use the link command with the symlink flag '-s': `ln -s /var/log/reactor.log /home/player/grid/reactor_link`.",
				ValidationType: "file_exists",
				Target:         "/home/player/grid/reactor_link",
			},
			{
				ID:             4,
				Title:          "Create Hard Link",
				Description:    "Create a hard link named `/home/player/grid/system_hard.log` that points to `/home/player/grid/system.log`.",
				Hint:           "Use the link command without any flags: `ln /home/player/grid/system.log /home/player/grid/system_hard.log`.",
				ValidationType: "file_exists",
				Target:         "/home/player/grid/system_hard.log",
			},
			{
				ID:             5,
				Title:          "Identify Archive Command",
				Description:    "What classic Linux utility is used to collect and combine multiple files into a single `.tar` archive or compress them as a `.tar.gz`?",
				Hint:           "This command is named after 'tape archive'. Type 'tar'.",
				ValidationType: "input_match",
				Expected:       "tar",
			},
			{
				ID:             6,
				Title:          "Magical Archives Packaging",
				Description:    "Create a compressed tar archive of your grid folder. Run `tar -czf /home/player/grid_archive.tar.gz -C /home/player grid` to package all your spell work.",
				Hint:           "Run the command: `tar -czf /home/player/grid_archive.tar.gz -C /home/player grid`.",
				ValidationType: "file_exists",
				Target:         "/home/player/grid_archive.tar.gz",
			},
		},
	},
	{
		ID:          8,
		Title:       "Chapter 8: The Guard Wards (Users & Group Security)",
		Description: "Understand authentication controls. Create new user profiles, establish magical groups, and shift file ownership.",
		Tasks: []Task{
			{
				ID:             1,
				Title:          "Summon User Profile",
				Description:    "Create a new user profile named `neville` on the server using sudo privileges.",
				Hint:           "Use the useradd command: `sudo useradd neville`.",
				ValidationType: "user_exists",
				Target:         "neville",
			},
			{
				ID:             2,
				Title:          "Assemble Magical Group",
				Description:    "Create a new security group named `gryffindor` on the server.",
				Hint:           "Use the groupadd command: `sudo groupadd gryffindor`.",
				ValidationType: "group_exists",
				Target:         "gryffindor",
			},
			{
				ID:             3,
				Title:          "Shift File Owner",
				Description:    "Change the owner of `/home/player/grid/hack.sh` to the newly created user `neville` using the change owner command.",
				Hint:           "Use chown with sudo: `sudo chown neville /home/player/grid/hack.sh`.",
				ValidationType: "owner_check",
				Target:         "/home/player/grid/hack.sh",
				Expected:       "neville",
			},
			{
				ID:             4,
				Title:          "Shift File Group",
				Description:    "Change the group ownership of `/home/player/grid/hack.sh` to the group `gryffindor`.",
				Hint:           "Use chgrp or chown with sudo: `sudo chgrp gryffindor /home/player/grid/hack.sh` or `sudo chown :gryffindor /home/player/grid/hack.sh`.",
				ValidationType: "group_check",
				Target:         "/home/player/grid/hack.sh",
				Expected:       "gryffindor",
			},
			{
				ID:             5,
				Title:          "Home Directories Parent",
				Description:    "In which absolute directory path are the personal home folders of non-root users created by default in Linux?",
				Hint:           "This is the top-level directory containing folders like 'player' or 'neville'. Type '/home'.",
				ValidationType: "input_match",
				Expected:       "/home",
			},
			{
				ID:             6,
				Title:          "Active Group Memberships",
				Description:    "Run the command `groups` in the terminal to inspect the groups your active user belongs to, and enter the first group name returned in the output.",
				Hint:           "Type `groups` in the terminal. The output lists your groups; submit the first name (usually 'player').",
				ValidationType: "input_match",
				Expected:       "player",
			},
		},
	},
	{
		ID:          9,
		Title:       "Chapter 9: Advanced Spell Filters (Text Transmutation)",
		Description: "Master advanced file scanning. Slice specific columns, print line ranges, and sort records.",
		Tasks: []Task{
			{
				ID:             1,
				Title:          "Identify Log Head",
				Description:    "What command line utility is used to print only the first few lines (defaulting to 10) of a text file?",
				Hint:           "This command is the opposite of 'tail'. Type 'head'.",
				ValidationType: "input_match",
				Expected:       "head",
			},
			{
				ID:             2,
				Title:          "Identify Log Tail",
				Description:    "What command line utility is used to output the last few lines of a file or follow real-time additions with `-f`?",
				Hint:           "This command is named after the back portion of an animal. Type 'tail'.",
				ValidationType: "input_match",
				Expected:       "tail",
			},
			{
				ID:             3,
				Title:          "Slice Column Fields",
				Description:    "Which command line utility is used to cut out and display selected columns or character fields from lines of files?",
				Hint:           "This shares a name with cutting paper. Type 'cut'.",
				ValidationType: "input_match",
				Expected:       "cut",
			},
			{
				ID:             4,
				Title:          "Sort Scroll Lines",
				Description:    "What command sorts lines of text files alphabetically or numerically?",
				Hint:           "Type the verb that means organizing things in order. Type 'sort'.",
				ValidationType: "input_match",
				Expected:       "sort",
			},
			{
				ID:             5,
				Title:          "Remove Duplicate Entries",
				Description:    "What command filters adjacent duplicate lines in a sorted file, showing unique lines?",
				Hint:           "This stands for 'uniq'. Type 'uniq'.",
				ValidationType: "input_match",
				Expected:       "uniq",
			},
			{
				ID:             6,
				Title:          "Word Count Spell",
				Description:    "What is the Linux command used to count the number of lines, words, and characters in a file?",
				Hint:           "Type 'wc' (which stands for Word Count).",
				ValidationType: "input_match",
				Expected:       "wc",
			},
		},
	},
	{
		ID:          10,
		Title:       "Chapter 10: The Chamber of Shells (Advanced Vault)",
		Description: "You have breached the Chamber of Shells! You must locate the chamber core logs, extract the sleeping password, bypass the basilisk's cooling ducts, and master automation cron daemons.",
		Tasks: []Task{
			{
				ID:             1,
				Title:          "Locate Chamber Logs",
				Description:    "Search inside `/var` for any files with 'reactor' in their name by running `find /var -name '*reactor*'`.",
				Hint:           "The find command displays the full path of matching files. Type the absolute file path you found.",
				ValidationType: "input_match",
				Target:         "reactor_log_path",
				Expected:       "/var/log/reactor.log",
			},
			{
				ID:             2,
				Title:          "Extracting the Sleep password",
				Description:    "Search `/var/log/reactor.log` using `grep` for the word 'CRITICAL' to extract the basilisk's containment sleeping password.",
				Hint:           "Run 'grep CRITICAL /var/log/reactor.log'. Look for the PASSWORD code printed.",
				ValidationType: "input_match",
				Target:         "reactor_password",
				Expected:       "CORE_TEMP_CRITICAL_9982",
			},
			{
				ID:             3,
				Title:          "Lock the Cooling Valves",
				Description:    "Filter all lines in `/var/log/reactor.log` containing 'ACTIVE' and redirect them to a new file named `/home/player/grid/active_sectors.txt` to seal the chamber.",
				Hint:           "Use grep and redirection: 'grep ACTIVE /var/log/reactor.log > /home/player/grid/active_sectors.txt'.",
				ValidationType: "file_contains",
				Target:         "/home/player/grid/active_sectors.txt",
				Expected:       "ACTIVE",
			},
			{
				ID:             4,
				Title:          "Magical Time Daemon",
				Description:    "What system daemon utility executes scheduled commands or scripts at specified recurring times in Linux?",
				Hint:           "This service uses files named 'crontab'. Type 'cron' or 'crond'.",
				ValidationType: "input_match",
				Expected:       "cron",
			},
			{
				ID:             5,
				Title:          "Recursive File Finder",
				Description:    "What is the name of the tool used to search for files and folders recursively based on attributes like name, size, or type?",
				Hint:           "This shares a name with the verb meaning to search for something. Type 'find'.",
				ValidationType: "input_match",
				Expected:       "find",
			},
			{
				ID:             6,
				Title:          "Locate Files by User",
				Description:    "Find all files in `/home/player/grid` owned by user `neville`. Run `find /home/player/grid -user neville` in the terminal. Enter the absolute path of the file that matches.",
				Hint:           "Run `find /home/player/grid -user neville`. Submit the file path returned: `/home/player/grid/hack.sh`.",
				ValidationType: "input_match",
				Expected:       "/home/player/grid/hack.sh",
			},
		},
	},
	{
		ID:          11,
		Title:       "Chapter 11: The Environment Sanctum (Environments & PATH)",
		Description: "Learn how the shell stores variables, customize your environment profiles, and configure command aliases.",
		Tasks: []Task{
			{
				ID:             1,
				Title:          "Sensing Environment Variables",
				Description:    "What command displays all active environment variables in your current shell session?",
				Hint:           "Type 'env' or 'printenv'.",
				ValidationType: "input_match",
				Expected:       "env",
			},
			{
				ID:             2,
				Title:          "Conjuring the PATH Rune",
				Description:    "What environment variable holds the colon-separated list of directories where the system looks for executable commands?",
				Hint:           "Type 'PATH' (case-sensitive, all caps).",
				ValidationType: "input_match",
				Expected:       "PATH",
			},
			{
				ID:             3,
				Title:          "Creating a Secret Variable",
				Description:    "Export an environment variable named `SECRET_SPELL` with the value `ExpectoPatronum` in your terminal: `export SECRET_SPELL=\"ExpectoPatronum\"`. Then, redirect the value of this variable into a new file named `/home/player/grid/spell_env.txt` by running: `echo $SECRET_SPELL > /home/player/grid/spell_env.txt`.",
				Hint:           "Type: `export SECRET_SPELL=\"ExpectoPatronum\"` and then `echo $SECRET_SPELL > /home/player/grid/spell_env.txt`.",
				ValidationType: "file_contains",
				Target:         "/home/player/grid/spell_env.txt",
				Expected:       "ExpectoPatronum",
			},
			{
				ID:             4,
				Title:          "Sensing the Default Shell",
				Description:    "What environment variable contains the path to the current user's default command shell (e.g. `/bin/bash`)?",
				Hint:           "Type 'SHELL' (case-sensitive, all caps).",
				ValidationType: "input_match",
				Expected:       "SHELL",
			},
			{
				ID:             5,
				Title:          "Defining a Spell Alias",
				Description:    "Write an alias named `reveal` that maps to `ls -la`. Put this alias in the user's configuration file `/home/player/.bashrc` by appending `alias reveal=\"ls -la\"` to the end of the file.",
				Hint:           "Run: `echo 'alias reveal=\"ls -la\"' >> /home/player/.bashrc` in your terminal.",
				ValidationType: "file_contains",
				Target:         "/home/player/.bashrc",
				Expected:       "alias reveal=",
			},
			{
				ID:             6,
				Title:          "Locating Executables",
				Description:    "What command is used to locate the binary executable file path for a given command?",
				Hint:           "Type 'which'. E.g. `which ls` returns `/usr/bin/ls`.",
				ValidationType: "input_match",
				Expected:       "which",
			},
		},
	},
	{
		ID:          12,
		Title:       "Chapter 12: Service Control & Daemons (Systemctl & Logs)",
		Description: "Understand background system engines. Manage essential daemons and query log databases.",
		Tasks: []Task{
			{
				ID:             1,
				Title:          "The System Engine Daemon",
				Description:    "What is the default init system and system manager used in modern Ubuntu/Linux systems to boot the system and manage services (PID 1)?",
				Hint:           "Type 'systemd'.",
				ValidationType: "input_match",
				Expected:       "systemd",
			},
			{
				ID:             2,
				Title:          "Managing Active Daemons",
				Description:    "What command-line tool is used to inspect and control the systemd system and service manager?",
				Hint:           "Type 'systemctl'.",
				ValidationType: "input_match",
				Expected:       "systemctl",
			},
			{
				ID:             3,
				Title:          "Inspecting Service Status",
				Description:    "What systemctl subcommand is used to check the detailed status of a service (e.g., systemctl ____ cron)?",
				Hint:           "Type 'status'.",
				ValidationType: "input_match",
				Expected:       "status",
			},
			{
				ID:             4,
				Title:          "Magical Log Journal",
				Description:    "What systemd utility is used to query and view logs generated by systemd's logging service?",
				Hint:           "Type 'journalctl'.",
				ValidationType: "input_match",
				Expected:       "journalctl",
			},
			{
				ID:             5,
				Title:          "Stopping Rogue Services",
				Description:    "If a systemd service gets stuck, what systemctl subcommand is used to immediately terminate and stop the service?",
				Hint:           "Type 'stop'.",
				ValidationType: "input_match",
				Expected:       "stop",
			},
			{
				ID:             6,
				Title:          "System Uptime Check",
				Description:    "What simple command displays how long the system has been running, along with the current time and load averages?",
				Hint:           "Type 'uptime'.",
				ValidationType: "input_match",
				Expected:       "uptime",
			},
		},
	},
	{
		ID:          13,
		Title:       "Chapter 13: Network Gateways & Portals (Netstat, SS, DNS)",
		Description: "Perform network diagnostics. Inspect listening port portals, resolve domains, and inspect device mappings.",
		Tasks: []Task{
			{
				ID:             1,
				Title:          "Inspecting Socket Connections",
				Description:    "What modern command is used to dump socket statistics, replacing the older `netstat` tool?",
				Hint:           "Type 'ss' (socket statistics).",
				ValidationType: "input_match",
				Expected:       "ss",
			},
			{
				ID:             2,
				Title:          "Resolving Domain Runes",
				Description:    "What command line tool is used to perform DNS lookup queries, returning IP addresses for hostnames?",
				Hint:           "Type 'nslookup'.",
				ValidationType: "input_match",
				Expected:       "nslookup",
			},
			{
				ID:             3,
				Title:          "Tracing Network Portals",
				Description:    "What tool prints the route (the path of routers) that packets take to reach a network host?",
				Hint:           "Type 'traceroute'.",
				ValidationType: "input_match",
				Expected:       "traceroute",
			},
			{
				ID:             4,
				Title:          "Portal IP Address Configuration",
				Description:    "What modern command is used to show or manipulate network routing, network devices, interfaces, and tunnels?",
				Hint:           "Type 'ip'. E.g. `ip addr` list interfaces.",
				ValidationType: "input_match",
				Expected:       "ip",
			},
			{
				ID:             5,
				Title:          "The Magical Hosts Rune",
				Description:    "In which file path is the local static table of domain names mapping directly to IP addresses located?",
				Hint:           "Type '/etc/hosts'.",
				ValidationType: "input_match",
				Expected:       "/etc/hosts",
			},
			{
				ID:             6,
				Title:          "Retrieving Remote Scrolls",
				Description:    "What classic command line utility is used to download files from web servers, supporting recursive background downloading?",
				Hint:           "Type 'wget'.",
				ValidationType: "input_match",
				Expected:       "wget",
			},
		},
	},
	{
		ID:          14,
		Title:       "Chapter 14: Stream Alchemy (Sed & Awk Filters)",
		Description: "Mutate inputs dynamically. Replace spell components via stream editing and slice records using layout delimiters.",
		Tasks: []Task{
			{
				ID:             1,
				Title:          "Stream Editor Transmutation",
				Description:    "What stream editor utility is used to perform basic text transformations on an input stream or file? (Hint: 3 letters).",
				Hint:           "Type 'sed'.",
				ValidationType: "input_match",
				Expected:       "sed",
			},
			{
				ID:             2,
				Title:          "Pattern Scanning Alchemy",
				Description:    "What powerful pattern scanning and processing language utility is named after its creators Aho, Weinberger, and Kernighan?",
				Hint:           "Type 'awk'.",
				ValidationType: "input_match",
				Expected:       "awk",
			},
			{
				ID:             3,
				Title:          "Cast Substitutions",
				Description:    "Use `sed` to substitute the word 'Voldemort' with 'He-Who-Must-Not-Be-Named' inside a new file. First, create `/home/player/grid/dark_magic.txt` containing 'Voldemort'. Then run sed to replace it and output the result to `/home/player/grid/cleansed_magic.txt`.",
				Hint:           "First create: `echo Voldemort > /home/player/grid/dark_magic.txt`. Then run: `sed 's/Voldemort/He-Who-Must-Not-Be-Named/g' /home/player/grid/dark_magic.txt > /home/player/grid/cleansed_magic.txt`.",
				ValidationType: "file_contains",
				Target:         "/home/player/grid/cleansed_magic.txt",
				Expected:       "He-Who-Must-Not-Be-Named",
			},
			{
				ID:             4,
				Title:          "Extracting Specific Fields",
				Description:    "Use `awk` to print the first field (column) of `/etc/passwd` using `:` as field separator, and write the output of the first 5 lines to `/home/player/grid/users_list.txt`.",
				Hint:           "Run: `awk -F: '{print $1}' /etc/passwd | head -n 5 > /home/player/grid/users_list.txt`.",
				ValidationType: "file_exists",
				Target:         "/home/player/grid/users_list.txt",
			},
			{
				ID:             5,
				Title:          "Regular Expression Sensing",
				Description:    "What symbol in regular expressions is used to match the beginning of a line?",
				Hint:           "Type the caret symbol '^'.",
				ValidationType: "input_match",
				Expected:       "^",
			},
			{
				ID:             6,
				Title:          "Regular Expression Endings",
				Description:    "What symbol in regular expressions is used to match the end of a line?",
				Hint:           "Type the dollar sign symbol '$'.",
				ValidationType: "input_match",
				Expected:       "$",
			},
		},
	},
	{
		ID:          15,
		Title:       "Chapter 15: The Sorcerer's Gate (SSH & Secure Transports)",
		Description: "You have reached the final gate of Codewarts! Master remote terminal portals, key pair authorization, and secure incremental transports.",
		Tasks: []Task{
			{
				ID:             1,
				Title:          "Establishing the Secure Portal",
				Description:    "What command is used to log into remote machines and execute commands securely over encrypted connections?",
				Hint:           "Type 'ssh'.",
				ValidationType: "input_match",
				Expected:       "ssh",
			},
			{
				ID:             2,
				Title:          "Magical Key Generator",
				Description:    "What tool is used to generate authentication key pairs for SSH secure logins?",
				Hint:           "Type 'ssh-keygen'.",
				ValidationType: "input_match",
				Expected:       "ssh-keygen",
			},
			{
				ID:             3,
				Title:          "Secure Scroll Copy",
				Description:    "What command-line utility copies files between hosts on a network securely over an SSH connection?",
				Hint:           "Type 'scp'.",
				ValidationType: "input_match",
				Expected:       "scp",
			},
			{
				ID:             4,
				Title:          "Incremental Synchronization Ritual",
				Description:    "What fast and extraordinarily versatile file copying tool is used to sync files locally and remotely by transferring only differences?",
				Hint:           "Type 'rsync'.",
				ValidationType: "input_match",
				Expected:       "rsync",
			},
			{
				ID:             5,
				Title:          "Default Portal Port",
				Description:    "What port number does the SSH service listen on by default?",
				Hint:           "Type '22'.",
				ValidationType: "input_match",
				Expected:       "22",
			},
			{
				ID:             6,
				Title:          "Checking SSH Key Location",
				Description:    "In which hidden subdirectory of the user's home directory are default SSH private and public keys stored?",
				Hint:           "Type '/home/player/.ssh'.",
				ValidationType: "input_match",
				Expected:       "/home/player/.ssh",
			},
		},
	},
}

// VerifyTask verifies a task completion inside a user's container.
func VerifyTask(containerID string, task Task, userInput string) (bool, error) {
	log.Printf("Verifying task %d (Type: %s, Target: %s, Input: '%s')",
		task.ID, task.ValidationType, task.Target, userInput)

	switch task.ValidationType {
	case "input_match":
		// Direct string match against expected value (case-insensitive and trimmed)
		cleanInput := strings.TrimSpace(strings.ToLower(userInput))
		cleanExpected := strings.TrimSpace(strings.ToLower(task.Expected))
		return cleanInput == cleanExpected, nil

	case "file_exists":
		// Check if file exists inside container
		cmd := []string{"test", "-f", task.Target}
		_, exitCode, err := container.RunExecCommand(containerID, "player", cmd)
		if err != nil {
			return false, err
		}
		return exitCode == 0, nil

	case "dir_exists":
		// Check if directory exists inside container
		cmd := []string{"test", "-d", task.Target}
		_, exitCode, err := container.RunExecCommand(containerID, "player", cmd)
		if err != nil {
			return false, err
		}
		return exitCode == 0, nil

	case "permission_check":
		// Check if file is executable inside container
		cmd := []string{"test", "-x", task.Target}
		_, exitCode, err := container.RunExecCommand(containerID, "player", cmd)
		if err != nil {
			return false, err
		}
		return exitCode == 0, nil

	case "file_contains":
		// Check if file exists and contains the expected string
		// First verify file exists
		existCmd := []string{"test", "-f", task.Target}
		_, exitCode, err := container.RunExecCommand(containerID, "player", existCmd)
		if err != nil || exitCode != 0 {
			return false, nil
		}

		// Read file contents
		catCmd := []string{"cat", task.Target}
		content, _, err := container.RunExecCommand(containerID, "player", catCmd)
		if err != nil {
			return false, err
		}

		cleanContent := strings.ToLower(content)
		cleanExpected := strings.ToLower(task.Expected)
		return strings.Contains(cleanContent, cleanExpected), nil

	case "process_running":
		// Check if a process is running inside the container
		cmd := []string{"pgrep", "-f", task.Target}
		_, exitCode, err := container.RunExecCommand(containerID, "player", cmd)
		if err != nil {
			return false, err
		}
		return exitCode == 0, nil

	case "process_not_running":
		// Check if a process is NOT running inside the container
		cmd := []string{"pgrep", "-f", task.Target}
		_, exitCode, err := container.RunExecCommand(containerID, "player", cmd)
		if err != nil {
			return false, err
		}
		return exitCode != 0, nil

	case "user_exists":
		// Check if user exists on the system
		cmd := []string{"id", "-u", task.Target}
		_, exitCode, err := container.RunExecCommand(containerID, "root", cmd)
		if err != nil {
			return false, err
		}
		return exitCode == 0, nil

	case "group_exists":
		// Check if group exists on the system
		cmd := []string{"getent", "group", task.Target}
		_, exitCode, err := container.RunExecCommand(containerID, "root", cmd)
		if err != nil {
			return false, err
		}
		return exitCode == 0, nil

	case "owner_check":
		// Check if user owner matches expected owner
		cmd := []string{"stat", "-c", "%U", task.Target}
		owner, _, err := container.RunExecCommand(containerID, "root", cmd)
		if err != nil {
			return false, err
		}
		return strings.TrimSpace(owner) == task.Expected, nil

	case "group_check":
		// Check if group matches expected group
		cmd := []string{"stat", "-c", "%G", task.Target}
		group, _, err := container.RunExecCommand(containerID, "root", cmd)
		if err != nil {
			return false, err
		}
		return strings.TrimSpace(group) == task.Expected, nil

	case "hostname_check":
		// Check if user input matches container hostname
		cmd := []string{"hostname"}
		hostname, _, err := container.RunExecCommand(containerID, "player", cmd)
		if err != nil {
			return false, err
		}
		return strings.TrimSpace(userInput) == strings.TrimSpace(hostname), nil

	default:
		return false, fmt.Errorf("unknown validation type: %s", task.ValidationType)
	}
}

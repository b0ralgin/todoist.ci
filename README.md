## Desciption 
This project is a TUI for todoist client (link). The project written in golang, uses gocui (link) as a graphical framework. 

#Installation 
## Build from source 
You need Go 1.15
```bash 
git clone $REPO_LINK
cd todoist-cui
make install 
```

## Packages 
 AUR[link]
 [comment]: DEB[Link]
 [comment]: OBS[Link]
[comment]: ### Mac OS 

# Main window 
Pic 
### Add/Edit Task
Pic 

## Keybindings 
|Key | Description |
|-----| ---------- |
|Ctrl-C | Exit program | 
|Ctrl-N | Create new task | 
|Ctrl-E | Edit current task |
|Arrows Up/Down | change current task | 
|Ctrl-P | Change task's project| 
|Tab | switch between fields | 
|Up/Down | change priority (!!!!- Highest priority, 0 - no Priority) |
|Up/Down | change Due to date | 
| Ctrl-_ | in Due to - erase Due to date |


## List of features 
 - [x] List of tasks 
 - [x] Add task
 - [x] Edit task
 - [x] Change task's project 
 - [x] Set task done 
 - [x] Delete task 
 - [] Custom keybindings 
 - [] User filter 
 - [] Task labels
 - [] Karma
 - [] Collaboration tools 

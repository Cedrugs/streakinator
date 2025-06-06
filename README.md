# streakinator

**streakinator** is a simple, interactive Go utility that creates backdated, empty Git commits to populate your GitHub contribution graph.


## ⚠️ Warning

This utility is provided strictly for **legitimate, emergency, or maintenance scenarios**.

Do **not** use it to:
- Fabricate GitHub activity
- Pad your résumé
- Trick interviewers or recruiters
- Fake a portfolio
- Commit academic fraud

> Misusing this tool is unethical and may harm your professional or academic reputation.  
By running streakinator, you acknowledge and accept full responsibility for its use.


## Features

- Interactive CLI wizard  
- Creates directory and `git init` if missing  
- Detects or creates branches  
- Fixed or randomized commits per day  
- Single date or date range  
- Backdated commits with `GIT_AUTHOR_DATE` and `GIT_COMMITTER_DATE`  
- Custom commit message support  
- Push reminder after commit creation  

## Prerequisites

- **Go 1.16+**
- **Git** installed
- OS permissions to write files and run git commands


## Installation

```bash
git clone https://github.com/yourusername/streakinator.git
cd streakinator
go build -o streakinator.exe ./cmd/streakinator
```

(Optional) move the binary to a folder in your `$PATH`:

```bash
mv streakinator /usr/local/bin/
```

Windows users: move `streakinator.exe` into a folder in your `%PATH%`.

## Usage

Run:

```bash
./streakinator
```

The tool will walk you through:
- Setting the repo directory
- Branch creation/selection
- Commit configuration
- Date selection
- Commit message

After it's done, you’ll be told how to push:

```bash
git push origin <branch>
```


## License

MIT License. Use responsibly and ethically.
# pulla

Backup personal GitHub and starred repositories to the local file system.

## Usage

Single run

```bash
pulla --dest "$PWD/repos" --token "github_pat_"
```

Pull new changes every 24 hours

```bash
pulla --dest "$PWD/repos" --token "github_pat_" --daemon --interval 24
```

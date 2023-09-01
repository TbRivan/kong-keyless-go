## Custom Plugin for Authentication - Keyless Auth

### How to build

- Install all dependencies

```bash
  go mod install
```

- If u wanna run it inside Docker Container, so u need to build it for Linux Environment, dont worry we already create the Makefile for every platform
- For Linux 64 bit

```bash
  make build-linux
```

- For Windows 64 bit

```bash
  make build-win
```

- For MacOS 64 bit Intel

```bash
  make build-mac
```

- For MacOS 64 bit Apple Silicon

```bash
  make build-mac-silicon
```

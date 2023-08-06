function watch-go-tests() {
  if [[ $# -eq 0 ]]; then
    echo "Usage: watch-go-tests <path to test files>"
    return 1
  fi

  local test_path="$1"
  fswatch -0 -o *.go | xargs -0 -n1 -I{} go test -v "$test_path"
}

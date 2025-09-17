#!/bin/bash

output_path="."
output_name="crawler"
joined="${output_path}/${output_name}"

# Parse options:
# -r skips program run
# -t skips tests
# -tr or -rt skips both
#
# The program is always built unless unit tests fail.
# So if you skip tests the program will be built regardless of them.
# Unless the build process itself fails of course.
# If the build fails and program run is not disabled it will be skipped,
# and you will get a message saying so.
run_tests=true
run_program=true

while getopts "tr" opt; do
  case $opt in
    r)
      run_program=false
      ;;
    t)
      run_tests=false
      ;;
    *)
      echo "Invalid option" >&2
      exit 1
      ;;
  esac
done

if $run_tests; then
  # Run Unit tests
  echo -e "----- Running unit tests -----" 
  test_output=$(go test ./... 2>&1)
  test_status=$?
  if [ $test_status -ne 0 ]; then
    echo -e "Some tests failed.\n"
    echo "$test_output"
    echo ""
    exit 1
  else
    echo "All is well."
  fi
fi

echo -e "\n----- Building into ${joined} -----" &&
go build -o $joined
build_status=$?
if [ $build_status -ne 0 ] && [ "$run_program" == true ]; then
  echo -e "Build failed, skipping program run.\n"
  run_program=false
  exit 1
fi

if $run_program; then
  echo -e "\n----- Running ${joined} -----\n" &&
  $joined
fi

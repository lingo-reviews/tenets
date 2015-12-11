Generating Python API
=====================

Installing Under Ubuntu
-----------------------

    mkvirtualenv --no-site-packages grpc
    git clone https://github.com/google/protobuf.git
    git clone https://github.com/grpc/grpc.git

For protobuf, just following the README instructions.

For grpc, I had to edit the Makefile to remove csharp targets.

    make
    sudo make install
    pip install tox
    tools/run_tests/build_python.sh 2.7



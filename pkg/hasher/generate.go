package hasher

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i PasswordHasher -o ./mocks/ -s "_minimock.go"

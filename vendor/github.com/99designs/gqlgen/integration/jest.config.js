module.exports = {
    testEnvironment: "node",
    testMatch: ["<rootDir>/**/*-test.js"],
    testPathIgnorePatterns: ["<rootDir>/node_modules/"],
    moduleFileExtensions: ["js"],
    modulePaths: ["<rootDir>/node_modules"],
    transform: {
        '^.+\\.jsx?$': 'babel-jest',
    },
};

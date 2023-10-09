module.exports = {
  env: {
    browser: true,
    es6: true
  },
  extends: ["eslint:recommended", "prettier"],
  plugins: ["prettier"],
  globals: {
    Atomics: "readonly",
    SharedArrayBuffer: "readonly"
  },
  overrides: [
    {
      files: ["**/*.ts", "**/*.tsx"],
      extends: [
        "airbnb-base",
        "airbnb-typescript",
        "prettier",
        "plugin:@typescript-eslint/recommended",
        "plugin:import/typescript"
      ],
      parser: "@typescript-eslint/parser",
      parserOptions: {
        ecmaFeatures: { jsx: true },
        project: "./tsconfig.eslint.json",
        tsconfigRootDir: __dirname
      },
      plugins: ["@typescript-eslint", "prettier"],
      rules: {
        "no-continue": "off",
        "no-multi-assign": "off",
        "no-param-reassign": "off",
        "no-plusplus": "off",
        "no-return-assign": "off"
      }
    }
  ]
};

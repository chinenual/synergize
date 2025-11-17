import js from "@eslint/js";
import globals from "globals";
import { defineConfig, globalIgnores } from "eslint/config";
/*eslint no-unused-vars: ["error", { "args": "none" }]*/

export default defineConfig([
    globalIgnores([
	"bindings/**",
	"dist/**",
	"public/**",
	"runtime-debug.js",
	"runtime.js"
    ]),
    {
	files: ["**/*.{js,mjs,cjs}"],
	
	plugins: { js },
	extends: [
	    "js/recommended",
	],
	languageOptions: {
	    globals: globals.browser
	},
	"rules": {
            "no-unused-vars": ["error", {
		"args": "none",
            }]
	},
    },
]);

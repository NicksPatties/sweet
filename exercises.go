package main

type exercise struct {
	name string
	text string
}

func getSampleExercises() []exercise {
	return []exercise{
		{
			"hello.go",
			"package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}",
		},
		{
			"hello.java",
			"class HelloWorld {\n\tpublic static void main(String[] args) {\n\t\tSystem.out.println(\"Hello, World!\");\n\t}\n}",
		},
		{
			"hello.js",
			"function sayHello() {\n  console.log('Hello, World!');\n}\n\nsayHello();",
		},
		{
			"hello.py",
			"def sayHello():\n\tprint('Hello, World!')\n\nif __name__ == '__main__':\n\tsayHello()",
		},
		{
			"hello.sql",
			"CREATE TABLE helloworld (phrase TEXT);\nINSERT INTO helloworld VALUES (\"Hello, World!\");\nINSERT INTO helloworld VALUES (\"Goodbye, World!\");\nSELECT COUNT(*) FROM helloworld;\n",
		},
	}
}

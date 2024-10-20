let fileContent = Tibeb.readFile(`${__dirname}/test.txt`);
const writeFIle = Tibeb.WriteFile(
  `${__dirname}/test.txt`,
  "this is updated text"
);

console.log(fileContent);

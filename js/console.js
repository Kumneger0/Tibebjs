// Test console.log
console.log("Testing console.log:", 42, { name: "Alice" });

// Test console.info
console.info("This is an informational message");

// Test console.warn
console.warn("Warning: Something might be wrong");

// Test console.error
console.error("Error occurred:", new Error("Test error"));

// Test console.debug
console.debug("Debugging information");

// Test console.assert
console.assert(false, "This assertion should fail");
console.assert(true, "This assertion should not print");

// Test console.clear
console.clear();

// Test console.count
console.count("counter1");
console.count("counter1");
console.count("counter2");
console.count("counter1");

// Test console.countReset
console.countReset("counter1");
console.count("counter1");

// Test console.group and console.groupEnd
console.group("Group 1");
console.log("Inside group 1");
console.group("Nested group");
console.log("Inside nested group");
console.groupEnd();
console.groupEnd();

// Test console.table
console.table([
  { name: "Alice", age: 30 },
  { name: "Bob", age: 25 },
]);

// Test console.time and console.timeEnd
console.time("timer1");
// Simulate some work
for (let i = 0; i < 1000000; i++) {}
console.timeEnd("timer1");

// Test console.trace
function functionC() {
  console.trace("Trace from functionC");
}

function functionB() {
  functionC();
}

function functionA() {
  functionB();
}

functionA();

console.log("All console methods have been tested.");

export const name = "name";

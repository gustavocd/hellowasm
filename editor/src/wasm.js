const go = new window.Go()

async function runGoWasm(rawData) {
  const result = await WebAssembly.instantiate(rawData, go.importObject)
  let oldLog = console.log;
  let stdOut = [];
  console.log = line => { stdOut.push(line) }
  await go.run(result.instance)
  console.log = oldLog
  return stdOut
}

export default runGoWasm

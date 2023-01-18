// let 1671110774101= {
//     "code\":100,\n    \"message\":\"记录不存在\"\n}\n\n\nsetValueByKey(ctx,\"$21330edbcb3242c1aacddfb69860d3fc.$resp.export \",1671110774101)"}

let a = false

async function sleep(time) {
    return new Promise((resolve) => {
        global.setTimeout()
        global.setTimeout(() => {
            resolve()
        }, time)

    })
}

async function f() {
    //等待
    await sleep(2000);
    a = true
}
f().then(r => {
    console.log(a);
    return a;
})
// f()

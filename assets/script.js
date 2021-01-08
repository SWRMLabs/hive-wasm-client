const go = new Go();
WebAssembly.instantiateStreaming(
    fetch("hive.wasm"),
    go.importObject
).then((result) => {
    go.run(result.instance);
    // MyFunc();
});
async function id() {
    try {
        let data = JSON.parse(await id());
         data = JSON.parse(data.val)
        // console.log(data.data);

        // console.log(await id())
        populateTable(data.data);
        
    } catch (err) {
        console.error("Caught exception", err);
    }
}





//selectors
const list  = document.getElementById("list");
const dashboardHeading = document.getElementById("dashboard_heading");
const TableData = document.getElementById("table_data");


//listners
list.addEventListener("click", e =>{
    e.preventDefault();
 const target = e.target.innerHTML;
 dashboardHeading.innerHTML = target;

switch(target){
    case 0: 
    target==="ID"
    id()
    break;
    case 1:
    target==="Status"
    //fetch status

    default:
    //default
}


})


// {/* <tr>
// <th scope="row">ID</th>
// <td>12D3KooWBQrdaiY1VYh9cdPyPkJ627deQ5shE4rx8VX9Lp2fdX5n</td>

// </tr>
// <tr>
// <th scope="row">Public Key</th>
// <td>12D3KooWBQrdaiY1VYh9cdPyPkJ627deQ5shE4rx8VX9Lp2fdX5n</td>
// </tr>

// <tr>
// <th scope="row">Adresses</th>
// <td>12D3KooWBQrdaiY1VYh9cdPyPkJ627deQ5shE4rx8VX9Lp2fdX5n
//    CAESIBexngmQaajsIsRRDp1CwmEGCt1rYVZq61zsIGhTlfgt
// </td>
// </tr> */}


function populateTable(data) {
    console.log(data)
    TableData.innerHTML = " ";
    for(const d in data){
        // console.log(`${d}: ${data[d]}`);
        const tr  = document.createElement('tr');
        const th = document.createElement('th');
        th.setAttribute('scope','row');
        th.innerHTML = `${d}`;
        tr.appendChild(th);
        const td = document.createElement('td');
        td.innerHTML = `${data[d]}`;
        tr.appendChild(td);

        TableData.appendChild(tr);

    }


}
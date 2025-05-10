
export default async function Profile() {

    const response = await fetch(`http://localhost:8080/api/profile`, {
        method: 'POST',
        credentials: 'include',  
        headers: {
          'Content-Type': 'application/json',
        },
        //here we should put the nickname of the owner of profile
        body: JSON.stringify({nickname: "iafriad"}),
    });
  
    if (!response.ok) {
        const resp = await response.text();
        console.log("Error get Data :", resp);
    } else {
        const data = await response.json();
        console.log("Success:");
        console.log(data);
    }
}
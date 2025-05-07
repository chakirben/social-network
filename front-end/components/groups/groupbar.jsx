"use client"
export default function Groupbar({onCreateGroup}) {
    return (
        <div className="groupbar">
            <div className="leftgroupbar">
                <img src="./images/arrow-left.svg" alt="Back" />
                <p>Groups</p>
            </div>
            <div id="btcreateGP" className="rightgroupbar" onClick={onCreateGroup}>
                <img className="imgplus" src="./images/plus.svg" alt="Add" />
                <button className="btcreategroup">Create Group</button>
            </div>
        </div>
    )
}
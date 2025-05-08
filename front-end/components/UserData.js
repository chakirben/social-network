export default function UserData({ usr }) {
    return (
      <div className="df">
        <img src={usr.image}/>
        <div className="df cl gp6">
          <h4>{usr.firstName} {usr.lastName}</h4>
          {usr.followers ? <h5>{usr.followers} followers</h5> : null}
        </div>
      </div>
    );
  }
  
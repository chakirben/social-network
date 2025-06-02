import { timePassed } from "@/public/utils/timePassed";
export default function DiscussionCard({ discussion }) {
  const avatarSrc =
    discussion.avatar && discussion.avatar.trim() !== ""
      ? discussion.avatar
      : "images/Avatars.png";

  return (
    <div className="discussionCard df spB center gp6">
      <div className="df  gp6 center spB">
        <img
          src={avatarSrc}
          alt={`${discussion.name} avatar`}
          className="avatar"
          style={{ width: 40, height: 40, borderRadius: "50%", marginRight: 10 }}
        />
        <div className="df cl gp6">
          <h4>{discussion.name}</h4>
          <p>{discussion.lastMessageContent || "No messages yet"}</p>
        </div>
      </div>
      <small>
        {discussion.lastMessageSentAt ? timePassed(discussion.lastMessageSentAt) : ""}
      </small>
    </div>
  );
}
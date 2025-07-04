import { timePassed } from "@/public/utils/timePassed";
import { useRouter } from "next/navigation";
import Avatar from "../avatar/avatar";
export default function DiscussionCard({ discussion }) {
  let router = useRouter()
  const avatarSrc =
    discussion.avatar && discussion.avatar.trim() !== ""
      ? discussion.avatar
      : "images/Avatars.png";

  return (
    <div className="discussionCard df spB center gp6"
      onClick={() => {
        const type = discussion.isGroup ? 'group' : 'user';
        const nameSlug = discussion.name.replace(/\s+/g, '_');
        router.push(`/chat/${type}${discussion.id}_${nameSlug}`);
      }}
    >
      <div className="df  gp12 center spB">
        <Avatar url={discussion.avatar} name={discussion.name} size={"big"} />
        <div className="df cl gp6">
          <h4>{discussion.name}</h4>
          <p>
            {discussion.lastMessageContent
              ? (discussion.lastMessageContent.length > 10
                ? discussion.lastMessageContent.slice(0, 10) + '...'
                : discussion.lastMessageContent)
              : 'No messages yet'}
          </p>
        </div>
      </div>
      <small>
        {discussion.lastMessageSentAt ? timePassed(discussion.lastMessageSentAt) : ""}
      </small>
    </div>
  );
}

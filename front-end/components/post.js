import { useState } from 'react';
import { useRouter } from 'next/navigation';
import ReactionGroup from './reactionGroup';
import { timePassed } from '@/public/utils/timePassed';

export default function Post({ pst }) {
    const router = useRouter();
    const [modalOpen, setModalOpen] = useState(false);
    const [modalImageSrc, setModalImageSrc] = useState('');

    const openImageModal = (src) => {
        setModalImageSrc(src);
        setModalOpen(true);
    };

    const closeModal = () => {
        setModalOpen(false);
        setModalImageSrc('');
    };

    return (
        <>
            <div className="post" onClick={() => router.push(`/post/${pst.id}`)}>
                <div className="content df cl gp12">
                    <div className="userData">
                        {pst?.avatar ? <img className="pic sm" src={`http://localhost:8080/${pst?.avatar}`} /> : ""}
                        <h4>{pst?.creator}</h4>
                        <h5>{"• " + timePassed(pst?.created_at)}</h5>
                    </div>
                    <div className="content">{pst?.content}</div>

                    {pst.image && (
                        <img
                            className="pic nrml"
                            src={`http://localhost:8080/${pst?.image}`}
                            onClick={(e) => {
                                e.stopPropagation(); // prevent post click navigation
                                openImageModal(`http://localhost:8080/${pst?.image}`);
                            }}
                            style={{ cursor: 'pointer' }}
                        />
                    )}
                </div>
                <ReactionGroup
                    className="reactionsContainer df gp12"
                    likeCount={pst.like_count}
                    dislikeCount={pst.dislike_count}
                    itemType="post"
                    itemId={pst.id}
                    userReaction={pst.user_reaction}/>
            </div>

            {modalOpen && (
                <div className="image-modal" onClick={closeModal}>
                    <img src={modalImageSrc} className="modal-img" />
                </div>
            )}
        </>
    );
}

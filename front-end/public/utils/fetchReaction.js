
export default async function fetchReaction(itemId, itemType, reactionValue) {
    console.log(itemId , itemType ,  reactionValue);
    
    try {
        const response = await fetch(`http://localhost:8080/api/reaction?itemType=${itemType}&itemId=${itemId}&reactionType=${reactionValue}`, {
            method: 'POST',
            credentials: "include"
        });

        if (!response.ok) {
            console.error("Server error:", response.statusText);
            return false;
        }

        // const data = await response.json();
        return true; // assumes API returns { success: true } on success
    } catch (error) {
        console.error("Network error:", error);
        return false;
    }
}

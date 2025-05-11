export default function User({ params }) {
    const { id } = params;

    return (
        <div>
            User ID: {id}
        </div>
    );
}
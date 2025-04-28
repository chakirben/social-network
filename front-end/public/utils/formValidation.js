export const validateDob = (dob) => {
    const age = (new Date() - new Date(dob)) / (1000 * 60 * 60 * 24 * 365);
    return age >= 16;
};

export const validateForm = (firstName, lastName, dob, email, password) => {
    if (!firstName || !lastName || !dob || !email || !password) {
        return { valid: false, message: "All required fields must be filled" };
    }

    const emailRegex = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$/;
    if (!emailRegex.test(email)) {
        return { valid: false, message: "Please enter a valid email address" };
    }

    if (password.length < 6) {
        return { valid: false, message: "Password must be at least 6 characters long" };
    }

    if (!validateDob(dob)) {
        return { valid: false, message: "User must be at least 16 years old" };
    }

    return { valid: true, message: "" };
};

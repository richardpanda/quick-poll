export const fetchGetPoll = async (id) => {
  try {
    const response = await fetch(`/v1/polls/${id}`);

    if (response.status === 500) {
      throw new Error('Something went wrong...');
    }

    const payload = await response.json();
    if (response.ok) {
      return payload;
    }
    throw new Error(payload.message);
  } catch (e) {
    throw e;
  }
};

export const fetchPostChoice = async (pollId, choiceId) => {
  try {
    const response = await fetch(`/v1/polls/${pollId}/choices/${choiceId}`, { method: 'POST' });

    if (response.status === 500) {
      throw new Error("Something went wrong...");
    }

    const payload = await response.json();
    if (response.ok) {
      return payload;
    }
    return new Error(payload.message);
  } catch (e) {
    throw e;
  }
};

export const fetchPostPoll = async ({ question, choices, checkIP: check_ip }) => {
  try {
    const opts = {
      method: 'POST',
      body: JSON.stringify({ question, choices, check_ip }),
    };
    const response = await fetch('/v1/polls', opts);

    if (response.status === 500) {
      throw new Error("Something went wrong...");
    }

    const payload = await response.json();
    if (response.ok) {
      return payload;
    }
    throw new Error(payload.message);
  } catch (e) {
    throw e;
  }
};

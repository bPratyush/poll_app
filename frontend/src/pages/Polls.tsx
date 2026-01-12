import { useState, useEffect, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { Poll } from '../types';
import { pollAPI } from '../services/api';
import { useAuth } from '../context/AuthContext';

// Polling interval in milliseconds (3 seconds)
const POLL_INTERVAL = 3000;

// Helper to check if user has seen the poll update
const getSeenUpdates = (): Record<number, string> => {
  try {
    return JSON.parse(localStorage.getItem('seenPollUpdates') || '{}');
  } catch {
    return {};
  }
};

const hasSeenUpdate = (pollId: number, updatedAt: string): boolean => {
  const seen = getSeenUpdates();
  return seen[pollId] === updatedAt;
};

function Polls() {
  const [polls, setPolls] = useState<Poll[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const { user } = useAuth();

  const fetchPolls = useCallback(async (showLoading = false) => {
    if (showLoading) setLoading(true);
    try {
      const response = await pollAPI.list();
      setPolls(response.data);
      setError('');
    } catch {
      setError('Failed to fetch polls');
    } finally {
      if (showLoading) setLoading(false);
    }
  }, []);

  // Initial fetch
  useEffect(() => {
    fetchPolls(true);
  }, [fetchPolls]);

  // Auto-refresh every 5 seconds
  useEffect(() => {
    const interval = setInterval(() => {
      fetchPolls(false);
    }, POLL_INTERVAL);

    return () => clearInterval(interval);
  }, [fetchPolls]);

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this poll?')) return;

    try {
      await pollAPI.delete(id);
      setPolls(polls.filter((p) => p.id !== id));
    } catch {
      setError('Failed to delete poll');
    }
  };

  if (loading) {
    return <div className="loading">Loading polls...</div>;
  }

  return (
    <div>
      <div className="polls-header">
        <h1>Active <span>Polls</span></h1>
        <Link to="/polls/new" className="btn btn-primary">
          Create Poll
        </Link>
      </div>

      {error && <div className="alert alert-error">{error}</div>}

      {polls.length === 0 ? (
        <div className="empty-state-hero">
          <h2>Welcome to PollApp</h2>
          <p className="empty-state-subtitle">Create interactive polls and gather opinions from your team, friends, or community.</p>
          <Link to="/polls/new" className="btn btn-primary btn-lg">Create Your First Poll</Link>
          
          <div className="empty-state-divider">
            <span>How it works</span>
          </div>
          
          <div className="empty-state-features">
            <div className="feature-item">
              <div className="feature-icon">1</div>
              <span>Create a poll with multiple options</span>
            </div>
            <div className="feature-item">
              <div className="feature-icon">2</div>
              <span>Share with others to collect votes</span>
            </div>
            <div className="feature-item">
              <div className="feature-icon">3</div>
              <span>View real-time results and voters</span>
            </div>
          </div>
        </div>
      ) : (
        <div className="polls-grid">
          {polls.map((poll) => (
            <div key={poll.id} className="poll-card">
              <div className="poll-card-header">
                <div>
                  <h3>{poll.title}</h3>
                  <span className="poll-card-meta">
                    by {poll.creator.username} • {poll.options.length} options
                  </span>
                </div>
                <div style={{ display: 'flex', gap: '0.5rem', flexWrap: 'wrap' }}>
                  {poll.poll_edited_after_vote && poll.user_voted_option_id && !hasSeenUpdate(poll.id, poll.updated_at) && (
                    <span className="poll-card-badge poll-card-badge-warning">
                      Updated
                    </span>
                  )}
                  {poll.user_voted_option_id && (
                    <span className="poll-card-badge">
                      Voted
                    </span>
                  )}
                </div>
              </div>
              {poll.description && <p>{poll.description}</p>}
              <div className="poll-card-actions">
                <Link to={`/polls/${poll.id}`} className="btn btn-primary">
                  View
                </Link>
                {user?.id === poll.creator.id && (
                  <>
                    <Link to={`/polls/${poll.id}/edit`} className="btn btn-secondary">
                      Edit
                    </Link>
                    <button onClick={() => handleDelete(poll.id)} className="btn btn-danger">
                      Delete
                    </button>
                  </>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default Polls;

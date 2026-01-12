import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Poll } from '../types';
import { pollAPI } from '../services/api';
import { useAuth } from '../context/AuthContext';

function Polls() {
  const [polls, setPolls] = useState<Poll[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const { user } = useAuth();

  useEffect(() => {
    fetchPolls();
  }, []);

  const fetchPolls = async () => {
    try {
      const response = await pollAPI.list();
      setPolls(response.data);
    } catch {
      setError('Failed to fetch polls');
    } finally {
      setLoading(false);
    }
  };

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
          ✨ Create Poll
        </Link>
      </div>

      {error && <div className="alert alert-error">{error}</div>}

      {polls.length === 0 ? (
        <div className="empty-state">
          <div className="empty-state-icon">📭</div>
          <h2>No polls yet</h2>
          <p>Create your first poll to get started!</p>
          <Link to="/polls/new" className="btn btn-primary">✨ Create Your First Poll</Link>
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
                {poll.user_voted_option_id && (
                  <span className="poll-card-badge">
                    ✓ Voted
                  </span>
                )}
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

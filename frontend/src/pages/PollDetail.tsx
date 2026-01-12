import { useState, useEffect, useCallback } from 'react';
import { useParams, Link } from 'react-router-dom';
import { Poll, User } from '../types';
import { pollAPI } from '../services/api';
import { useAuth } from '../context/AuthContext';

// Polling interval in milliseconds (5 seconds)
const POLL_INTERVAL = 5000;

// Mark poll update as seen in localStorage
const markUpdateAsSeen = (pollId: number, updatedAt: string) => {
  try {
    const seen = JSON.parse(localStorage.getItem('seenPollUpdates') || '{}');
    seen[pollId] = updatedAt;
    localStorage.setItem('seenPollUpdates', JSON.stringify(seen));
  } catch {
    // Ignore localStorage errors
  }
};

function PollDetail() {
  const { id } = useParams<{ id: string }>();
  const [poll, setPoll] = useState<Poll | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [selectedOption, setSelectedOption] = useState<number | null>(null);
  const [voting, setVoting] = useState(false);
  const [showVoters, setShowVoters] = useState<number | null>(null);
  const [voters, setVoters] = useState<User[]>([]);
  const [loadingVoters, setLoadingVoters] = useState(false);
  const [isChangingVote, setIsChangingVote] = useState(false);
  const [showUpdateNotice, setShowUpdateNotice] = useState(true);
  const { user } = useAuth();

  const fetchPoll = useCallback(async (showLoading = false) => {
    if (showLoading) setLoading(true);
    try {
      const response = await pollAPI.get(Number(id));
      setPoll(response.data);
      if (response.data.user_voted_option_id) {
        setSelectedOption(response.data.user_voted_option_id);
      }
      setError('');
    } catch {
      setError('Failed to fetch poll');
    } finally {
      if (showLoading) setLoading(false);
    }
  }, [id]);

  // Initial fetch
  useEffect(() => {
    fetchPoll(true);
  }, [fetchPoll]);

  // Auto-refresh every 5 seconds
  useEffect(() => {
    const interval = setInterval(() => {
      fetchPoll(false);
    }, POLL_INTERVAL);

    return () => clearInterval(interval);
  }, [fetchPoll]);

  // Mark update as seen when user views the poll
  useEffect(() => {
    if (poll && poll.poll_edited_after_vote && poll.user_voted_option_id) {
      markUpdateAsSeen(poll.id, poll.updated_at);
    }
  }, [poll]);

  const handleVote = async () => {
    if (!selectedOption || !poll) return;

    setVoting(true);
    try {
      const response = await pollAPI.vote(poll.id, selectedOption);
      setPoll(response.data);
      setIsChangingVote(false);
    } catch (err: unknown) {
      const error = err as { response?: { data?: { error?: string } } };
      setError(error.response?.data?.error || 'Failed to vote');
    } finally {
      setVoting(false);
    }
  };

  const handleChangeVote = () => {
    setIsChangingVote(true);
    setSelectedOption(null);
  };

  const handleCancelChange = () => {
    setIsChangingVote(false);
    if (poll?.user_voted_option_id) {
      setSelectedOption(poll.user_voted_option_id);
    }
  };

  const handleShowVoters = async (optionId: number) => {
    setShowVoters(optionId);
    setLoadingVoters(true);
    try {
      const response = await pollAPI.getVoters(optionId);
      setVoters(response.data);
    } catch {
      setVoters([]);
    } finally {
      setLoadingVoters(false);
    }
  };

  const getTotalVotes = () => {
    return poll?.options.reduce((sum, opt) => sum + opt.vote_count, 0) || 0;
  };

  const getVotePercentage = (count: number) => {
    const total = getTotalVotes();
    return total > 0 ? Math.round((count / total) * 100) : 0;
  };

  if (loading) {
    return <div className="loading">Loading poll...</div>;
  }

  if (!poll) {
    return <div className="empty-state"><h2>Poll not found</h2></div>;
  }

  const hasVoted = poll.user_voted_option_id !== undefined && poll.user_voted_option_id !== null;
  const isOwner = user?.id === poll.creator.id;
  const showVotingUI = !hasVoted || isChangingVote;

  return (
    <div className="poll-detail">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
        <div>
          <h1>{poll.title}</h1>
          <p className="poll-detail-meta">
            Created by {poll.creator.username} • {new Date(poll.created_at).toLocaleDateString()}
          </p>
        </div>
        {isOwner && (
          <Link to={`/polls/${poll.id}/edit`} className="btn btn-secondary">
            Edit Poll
          </Link>
        )}
      </div>

      {poll.description && (
        <p className="poll-detail-description">{poll.description}</p>
      )}

      {error && <div className="alert alert-error">{error}</div>}

      {/* Poll Edited After Vote Notification */}
      {poll.poll_edited_after_vote && hasVoted && !isChangingVote && showUpdateNotice && (
        <div className="alert alert-warning" style={{ marginBottom: '1.5rem', position: 'relative' }}>
          <div style={{ flex: 1 }}>
            <strong>Poll Updated:</strong> This poll was modified after you voted. 
            You may want to review the options and change your vote if needed.
          </div>
          <button 
            onClick={() => setShowUpdateNotice(false)}
            style={{ 
              background: 'none', 
              border: 'none', 
              cursor: 'pointer', 
              padding: '0.25rem',
              color: 'inherit',
              fontSize: '1.25rem',
              lineHeight: 1
            }}
            aria-label="Dismiss"
          >
            ×
          </button>
        </div>
      )}

      <div className="poll-options">
        <p className="poll-options-title">
          {showVotingUI ? 'Choose an option' : 'Results'}
          {isChangingVote && <span style={{ fontWeight: 'normal', fontSize: '0.9rem' }}> (changing vote)</span>}
        </p>
        {poll.options.map((option) => (
          <div
            key={option.id}
            className={`poll-option ${selectedOption === option.id ? 'selected' : ''} ${!showVotingUI ? 'voted' : ''} ${!showVotingUI && poll.user_voted_option_id === option.id ? 'user-choice' : ''}`}
            onClick={() => showVotingUI && setSelectedOption(option.id)}
          >
            {showVotingUI && (
              <div className="poll-option-radio" />
            )}
            <span className="poll-option-text">{option.text}</span>
            {!showVotingUI && (
              <>
                <span
                  className="poll-option-count"
                  onClick={(e) => {
                    e.stopPropagation();
                    handleShowVoters(option.id);
                  }}
                  title="Click to see voters"
                >
                  {option.vote_count} vote{option.vote_count !== 1 ? 's' : ''} ({getVotePercentage(option.vote_count)}%)
                </span>
              </>
            )}
            {!showVotingUI && (
              <div className="poll-option-bar">
                <div
                  className="poll-option-bar-fill"
                  style={{ width: `${getVotePercentage(option.vote_count)}%` }}
                />
              </div>
            )}
          </div>
        ))}
      </div>

      {showVotingUI && (
        <div style={{ display: 'flex', gap: '1rem', flexWrap: 'wrap' }}>
          <button
            className="btn btn-primary"
            style={{ flex: 1 }}
            onClick={handleVote}
            disabled={!selectedOption || voting}
          >
            {voting ? 'Submitting...' : isChangingVote ? 'Update Vote' : 'Submit Vote'}
          </button>
          {isChangingVote && (
            <button
              className="btn btn-secondary"
              onClick={handleCancelChange}
            >
              Cancel
            </button>
          )}
        </div>
      )}

      {hasVoted && !isChangingVote && (
        <div className="poll-footer">
          <p className="poll-total-votes">Total votes: {getTotalVotes()}</p>
          <p className="poll-vote-hint">Click on vote counts to see who voted</p>
          <button
            className="btn btn-secondary"
            style={{ marginTop: '1rem' }}
            onClick={handleChangeVote}
          >
            Change My Vote
          </button>
        </div>
      )}

      <div style={{ marginTop: '2rem', textAlign: 'center' }}>
        <Link to="/" className="btn btn-secondary">
          Back to Polls
        </Link>
      </div>

      {/* Voters Modal */}
      {showVoters !== null && (
        <div className="modal-overlay" onClick={() => setShowVoters(null)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h2>Voters for "{poll.options.find((o) => o.id === showVoters)?.text}"</h2>
              <button className="modal-close" onClick={() => setShowVoters(null)}>×</button>
            </div>
            {loadingVoters ? (
              <p>Loading voters...</p>
            ) : voters.length === 0 ? (
              <p className="text-muted text-center">No votes yet</p>
            ) : (
              <ul className="voters-list">
                {voters.map((voter) => (
                  <li key={voter.id}>
                    <div className="voter-avatar">{voter.username.charAt(0).toUpperCase()}</div>
                    <div className="voter-info">
                      <strong>{voter.username}</strong>
                      <span>{voter.email}</span>
                    </div>
                  </li>
                ))}
              </ul>
            )}
          </div>
        </div>
      )}
    </div>
  );
}

export default PollDetail;

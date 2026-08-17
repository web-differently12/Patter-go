import React, { useState } from 'react';

export interface TopUpModalProps {
  isOpen: boolean;
  onClose: () => void;
  currency: 'EUR' | 'USD';
  onConfirmTopUp: (amount: number, method: 'card' | 'sepa' | 'btc_lightning') => Promise<void>;
}

export const TopUpModal: React.FC<TopUpModalProps> = ({
  isOpen,
  onClose,
  currency,
  onConfirmTopUp,
}) => {
  const [selectedAmount, setAmount] = useState<number>(100);
  const [method, setMethod] = useState<'card' | 'sepa' | 'btc_lightning'>('card');
  const [loading, setLoading] = useState(false);

  if (!isOpen) return null;

  const symbol = currency === 'EUR' ? 'EUR' : 'USD';

  const handleTopUp = async () => {
    setLoading(true);
    try {
      await onConfirmTopUp(selectedAmount, method);
      onClose();
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm p-4 font-sans text-zinc-100">
      <div className="w-full max-w-md bg-[#09090b] border border-zinc-800/80 rounded-2xl shadow-2xl p-6 space-y-5">
        {/* Header */}
        <div className="flex justify-between items-center border-b border-zinc-800/80 pb-4">
          <div>
            <h2 className="text-base font-semibold text-zinc-100">Recharger le Wallet KallFlow</h2>
            <p className="text-xs text-zinc-400 mt-0.5">Paiement instantane securise par Hyperswitch.</p>
          </div>
          <button onClick={onClose} className="text-zinc-400 hover:text-zinc-100 text-sm p-1 rounded-lg">✕</button>
        </div>

        {/* Amount Selector */}
        <div className="space-y-2">
          <label className="text-xs font-medium text-zinc-400">Montant du rechargement</label>
          <div className="grid grid-cols-4 gap-2">
            {[50, 100, 250, 500].map((amt) => (
              <button
                key={amt}
                type="button"
                onClick={() => setAmount(amt)}
                className={`py-2 rounded-lg text-xs font-mono font-medium transition-colors ${
                  selectedAmount === amt
                    ? 'bg-zinc-100 text-zinc-950 font-semibold'
                    : 'bg-zinc-900 border border-zinc-800 text-zinc-300 hover:border-zinc-700'
                }`}
              >
                {amt} {symbol}
              </button>
            ))}
          </div>
        </div>

        {/* Payment Method Selector */}
        <div className="space-y-2">
          <label className="text-xs font-medium text-zinc-400">Methode de Paiement (Hyperswitch)</label>
          <div className="space-y-2">
            <button
              type="button"
              onClick={() => setMethod('card')}
              className={`w-full p-3 rounded-xl border text-left text-xs transition-colors flex justify-between items-center ${
                method === 'card' ? 'bg-zinc-900/80 border-zinc-500 text-zinc-100' : 'bg-zinc-950 border-zinc-800 text-zinc-400'
              }`}
            >
              <div className="flex items-center gap-2">
                <span className="font-mono text-zinc-[500]">[CB]</span>
                <span className="font-medium text-zinc-200">Carte Bancaire (Visa / Mastercard)</span>
              </div>
              <span className="text-[10px] text-emerald-400 font-mono">Instant</span>
            </button>

            <button
              type="button"
              onClick={() => setMethod('sepa')}
              className={`w-full p-3 rounded-xl border text-left text-xs transition-colors flex justify-between items-center ${
                method === 'sepa' ? 'bg-zinc-900/80 border-zinc-500 text-zinc-100' : 'bg-zinc-950 border-zinc-800 text-zinc-400'
              }`}
            >
              <div className="flex items-center gap-2">
                <span className="font-mono text-zinc-[500]">[SEPA]</span>
                <span className="font-medium text-zinc-200">Prelevement SEPA Direct</span>
              </div>
              <span className="text-[10px] text-zinc-500 font-mono">24h - 48h</span>
            </button>

            <button
              type="button"
              onClick={() => setMethod('btc_lightning')}
              className={`w-full p-3 rounded-xl border text-left text-xs transition-colors flex justify-between items-center ${
                method === 'btc_lightning' ? 'bg-zinc-900/80 border-zinc-500 text-zinc-100' : 'bg-zinc-950 border-zinc-800 text-zinc-400'
              }`}
            >
              <div className="flex items-center gap-2">
                <span className="font-mono text-zinc-[500]">[BTC]</span>
                <span className="font-medium text-zinc-200">Bitcoin Lightning (OpenNode / LND)</span>
              </div>
              <span className="text-[10px] text-amber-400 font-mono">0.1s Sats</span>
            </button>
          </div>
        </div>

        {/* Actions */}
        <div className="flex justify-end gap-2 pt-3 border-t border-zinc-800/80">
          <button
            type="button"
            onClick={onClose}
            className="px-4 py-2 rounded-lg text-xs text-zinc-400 hover:text-zinc-200 transition-colors"
          >
            Annuler
          </button>
          <button
            type="button"
            onClick={handleTopUp}
            disabled={loading}
            className="bg-white hover:bg-zinc-200 text-zinc-950 font-medium text-xs px-4 py-2 rounded-lg transition-colors shadow-sm disabled:opacity-50"
          >
            {loading ? 'Paiement...' : `Payer ${selectedAmount} ${symbol}`}
          </button>
        </div>
      </div>
    </div>
  );
};

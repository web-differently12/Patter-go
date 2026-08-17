import React, { useState } from 'react';

interface McpIntegration {
  id: string;
  name: string;
  category: 'CRM' | 'Workspace' | 'Support' | 'Database' | 'DevOps' | 'Payments';
  description: string;
  authType: 'OAuth2' | 'API Key';
  connected: boolean;
  discoveredTools: string[];
  transport: 'HTTP/SSE' | 'WebSocket';
}

export const McpServersPage: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'catalog' | 'connected' | 'custom'>('catalog');
  const [selectedCategory, setSelectedCategory] = useState<string>('All');
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [inspectToolsModal, setInspectToolsModal] = useState<McpIntegration | null>(null);

  // Pre-Built Catalog of MCP Servers available to Tenants
  const [integrations, setIntegrations] = useState<McpIntegration[]>([
    {
      id: 'mcp-hubspot',
      name: 'HubSpot CRM MCP Server',
      category: 'CRM',
      description: 'Synchronisation fiches contacts, opportunites commerciales, deals et creation de notes post-appel.',
      authType: 'OAuth2',
      connected: true,
      transport: 'HTTP/SSE',
      discoveredTools: [
        'hubspot_search_contacts',
        'hubspot_create_contact',
        'hubspot_update_deal_stage',
        'hubspot_log_call_engagement',
      ],
    },
    {
      id: 'mcp-salesforce',
      name: 'Salesforce Sales Cloud MCP',
      category: 'CRM',
      description: 'Lecture/Ecriture des Leads, Accounts, Opportunities, Tasks et objets personnalises Salesforce.',
      authType: 'OAuth2',
      connected: true,
      transport: 'HTTP/SSE',
      discoveredTools: [
        'salesforce_find_lead',
        'salesforce_create_lead',
        'salesforce_create_task',
        'salesforce_get_opportunity_by_id',
      ],
    },
    {
      id: 'mcp-google-workspace',
      name: 'Google Workspace Calendar & Gmail MCP',
      category: 'Workspace',
      description: 'Verification de disponibilites Google Calendar, reservation de creneaux et envoi de mails de confirmation.',
      authType: 'OAuth2',
      connected: true,
      transport: 'HTTP/SSE',
      discoveredTools: [
        'google_calendar_list_free_slots',
        'google_calendar_create_event',
        'gmail_send_followup_email',
      ],
    },
    {
      id: 'mcp-zendesk',
      name: 'Zendesk Support MCP Server',
      category: 'Support',
      description: "Creation automatique de tickets de support client, mise a jour des statuts et qualification par l'agent IA.",
      authType: 'OAuth2',
      connected: false,
      transport: 'HTTP/SSE',
      discoveredTools: [
        'zendesk_create_ticket',
        'zendesk_get_customer_tickets',
        'zendesk_update_ticket_status',
      ],
    },
    {
      id: 'mcp-slack',
      name: 'Slack Notification & Channel MCP',
      category: 'Workspace',
      description: "Diffusion d'alertes urgentes et de resumes d'appels dans des canaux Slack d'equipe.",
      authType: 'OAuth2',
      connected: false,
      transport: 'HTTP/SSE',
      discoveredTools: [
        'slack_send_channel_message',
        'slack_send_dm_notification',
      ],
    },
    {
      id: 'mcp-notion',
      name: 'Notion Knowledge Base MCP',
      category: 'Workspace',
      description: 'Recherche semantique RAG et lecture des bases de connaissances Notion en temps reel.',
      authType: 'OAuth2',
      connected: false,
      transport: 'HTTP/SSE',
      discoveredTools: [
        'notion_search_pages',
        'notion_append_block_children',
      ],
    },
    {
      id: 'mcp-postgres',
      name: 'PostgreSQL Vector & Relational MCP',
      category: 'Database',
      description: 'Execution de requetes SQL securisees sous PostgreSQL RLS et recherche de vecteurs pgvector.',
      authType: 'API Key',
      connected: true,
      transport: 'HTTP/SSE',
      discoveredTools: [
        'postgres_execute_read_query',
        'postgres_similarity_vector_search',
        'postgres_get_table_schema',
      ],
    },
    {
      id: 'mcp-stripe',
      name: 'Stripe Billing & Subscriptions MCP',
      category: 'Payments',
      description: "Verification des paiements, statut d'abonnement et generation de liens de paiement securises.",
      authType: 'API Key',
      connected: false,
      transport: 'HTTP/SSE',
      discoveredTools: [
        'stripe_get_customer_subscription',
        'stripe_create_payment_link',
        'stripe_refund_charge',
      ],
    },
    {
      id: 'mcp-jira',
      name: 'Jira Software & Service Desk MCP',
      category: 'DevOps',
      description: 'Creation de tickets de bugs et suivi des developpements directement depuis les conversations voix.',
      authType: 'OAuth2',
      connected: false,
      transport: 'HTTP/SSE',
      discoveredTools: [
        'jira_create_issue',
        'jira_get_issue_status',
        'jira_add_comment',
      ],
    },
    {
      id: 'mcp-linear',
      name: 'Linear Issue Tracking MCP',
      category: 'DevOps',
      description: "Gestion fluide des projets et tickets d'ingenierie Linear pendant les reunions.",
      authType: 'API Key',
      connected: false,
      transport: 'HTTP/SSE',
      discoveredTools: [
        'linear_create_issue',
        'linear_list_team_issues',
      ],
    },
  ]);

  const toggleConnection = (id: string) => {
    setIntegrations((prev) =>
      prev.map((item) =>
        item.id === id ? { ...item, connected: !item.connected } : item
      )
    );
  };

  const filteredIntegrations = integrations.filter((item) => {
    const matchesCategory = selectedCategory === 'All' || item.category === selectedCategory;
    const matchesSearch =
      item.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      item.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
      item.discoveredTools.some((tool) => tool.toLowerCase().includes(searchQuery.toLowerCase()));
    const matchesTab =
      activeTab === 'catalog'
        ? true
        : activeTab === 'connected'
        ? item.connected
        : false;
    return matchesCategory && matchesSearch && matchesTab;
  });

  return (
    <div className="space-y-6 text-slate-100 p-6 bg-[#08090a] min-h-screen font-sans">
      {/* Header */}
      <div className="border-b border-zinc-800 pb-4 flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <span className="text-xs font-mono font-bold text-cyan-400 bg-cyan-950/60 px-2.5 py-1 rounded-full border border-cyan-800/80">
            Model Context Protocol Gateway & Integration Hub
          </span>
          <h1 className="text-2xl font-bold tracking-tight text-white mt-1">Catalogue de Serveurs MCP & APIs Tierces</h1>
          <p className="text-xs text-zinc-400">
            Selectionnez et connectez les serveurs MCP pour alimenter vos agents vocaux et visio avec vos outils metiers (HubSpot, Salesforce, Google, Slack, etc.).
          </p>
        </div>

        <button className="bg-cyan-500 hover:bg-cyan-400 text-zinc-950 font-bold px-4 py-2 rounded-xl text-xs transition-colors shadow-sm self-start md:self-auto">
          + Connecter un Serveur MCP Custom (HTTP/SSE)
        </button>
      </div>

      {/* Navigation Tabs & Search Controls */}
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 bg-zinc-950 p-4 rounded-2xl border border-zinc-800">
        <div className="flex items-center gap-2 text-xs">
          <button
            onClick={() => setActiveTab('catalog')}
            className={`px-4 py-2 rounded-xl font-bold transition-all ${
              activeTab === 'catalog' ? 'bg-zinc-100 text-zinc-950' : 'text-zinc-400 hover:text-white bg-zinc-900'
            }`}
          >
            Catalogue Complet ({integrations.length})
          </button>
          <button
            onClick={() => setActiveTab('connected')}
            className={`px-4 py-2 rounded-xl font-bold transition-all ${
              activeTab === 'connected' ? 'bg-emerald-500 text-zinc-950' : 'text-zinc-400 hover:text-white bg-zinc-900'
            }`}
          >
            Connectes ({integrations.filter((i) => i.connected).length})
          </button>
        </div>

        {/* Search & Category Filter */}
        <div className="flex flex-wrap items-center gap-2 w-full md:w-auto text-xs">
          <input
            type="text"
            placeholder="Rechercher par nom ou outil MCP..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="bg-zinc-900 border border-zinc-800 rounded-xl px-3 py-1.5 text-zinc-200 placeholder-zinc-500 text-xs focus:outline-none focus:border-cyan-500 w-full sm:w-64"
          />

          <select
            value={selectedCategory}
            onChange={(e) => setSelectedCategory(e.target.value)}
            className="bg-zinc-900 border border-zinc-800 rounded-xl px-3 py-1.5 text-zinc-200 text-xs focus:outline-none focus:border-cyan-500"
          >
            <option value="All">Toutes Categories</option>
            <option value="CRM">CRM & Ventes</option>
            <option value="Workspace">Workspace & Comm</option>
            <option value="Support">Support Client</option>
            <option value="Database">Bases de Donnees</option>
            <option value="Payments">Paiements & Facturation</option>
            <option value="DevOps">DevOps & Projets</option>
          </select>
        </div>
      </div>

      {/* Catalog Cards Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {filteredIntegrations.map((item) => (
          <div
            key={item.id}
            className={`bg-zinc-950 border rounded-2xl p-5 shadow-xl flex flex-col justify-between space-y-4 transition-all ${
              item.connected ? 'border-emerald-500/60 shadow-emerald-950/20' : 'border-zinc-800 hover:border-zinc-700'
            }`}
          >
            <div className="space-y-3">
              <div className="flex justify-between items-start">
                <span className="text-[10px] font-mono font-bold px-2 py-0.5 rounded bg-zinc-900 text-cyan-400 border border-zinc-800">
                  {item.category}
                </span>
                <span
                  className={`text-[10px] font-mono font-bold px-2 py-0.5 rounded border ${
                    item.connected
                      ? 'bg-emerald-950 text-emerald-400 border-emerald-800'
                      : 'bg-zinc-900 text-zinc-500 border-zinc-800'
                  }`}
                >
                  {item.connected ? 'CONNECTE' : 'DISPONIBLE'}
                </span>
              </div>

              <h3 className="font-bold text-white text-base leading-snug">{item.name}</h3>
              <p className="text-xs text-zinc-400 leading-relaxed">{item.description}</p>

              {/* Tools Discovery Pill Count */}
              <div className="pt-2 border-t border-zinc-900 text-xs flex justify-between items-center font-mono">
                <span className="text-zinc-500 text-[11px]">{item.discoveredTools.length} Outils Decouverts</span>
                <button
                  onClick={() => setInspectToolsModal(item)}
                  className="text-cyan-400 hover:text-cyan-300 text-[11px] underline"
                >
                  Inspecter tools/list →
                </button>
              </div>
            </div>

            {/* Actions */}
            <div className="pt-2 flex items-center justify-between gap-2">
              <span className="text-[10px] font-mono text-zinc-500">Auth: {item.authType}</span>
              <button
                onClick={() => toggleConnection(item.id)}
                className={`px-4 py-2 rounded-xl font-bold text-xs transition-all shadow-sm ${
                  item.connected
                    ? 'bg-zinc-900 hover:bg-zinc-800 text-rose-400 border border-zinc-800'
                    : 'bg-cyan-500 hover:bg-cyan-400 text-zinc-950'
                }`}
              >
                {item.connected ? 'Deconnecter' : 'Connecter / Authoriser'}
              </button>
            </div>
          </div>
        ))}
      </div>

      {/* Tools Inspection Drawer/Modal */}
      {inspectToolsModal && (
        <div className="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center p-4 z-50">
          <div className="bg-zinc-950 border border-zinc-800 rounded-2xl p-6 max-w-2xl w-full space-y-4 shadow-2xl">
            <div className="flex justify-between items-center border-b border-zinc-800 pb-3">
              <div>
                <span className="text-[10px] font-mono font-bold text-cyan-400 bg-cyan-950 px-2 py-0.5 rounded border border-cyan-800">
                  {inspectToolsModal.category} MCP Server
                </span>
                <h3 className="font-bold text-white text-lg mt-1">{inspectToolsModal.name}</h3>
              </div>
              <button
                onClick={() => setInspectToolsModal(null)}
                className="text-zinc-400 hover:text-white font-bold text-sm bg-zinc-900 px-3 py-1 rounded-lg border border-zinc-800"
              >
                ✕ Fermer
              </button>
            </div>

            <p className="text-xs text-zinc-400">{inspectToolsModal.description}</p>

            <div className="space-y-2">
              <label className="text-xs font-bold text-zinc-300 block">
                Liste des Outils Decouverts via RPC `tools/list` ({inspectToolsModal.discoveredTools.length})
              </label>

              <div className="bg-zinc-900 border border-zinc-800 rounded-xl p-3 max-h-60 overflow-y-auto font-mono text-xs space-y-2">
                {inspectToolsModal.discoveredTools.map((tool) => (
                  <div key={tool} className="flex justify-between items-center p-2 bg-zinc-950 rounded border border-zinc-800/80">
                    <span className="text-cyan-300 font-bold">{tool}</span>
                    <span className="text-[10px] text-zinc-500">Exécutable par Agent IA</span>
                  </div>
                ))}
              </div>
            </div>

            <div className="flex justify-end pt-2 border-t border-zinc-900">
              <button
                onClick={() => setInspectToolsModal(null)}
                className="bg-cyan-500 hover:bg-cyan-400 text-zinc-950 font-bold px-4 py-2 rounded-xl text-xs"
              >
                Compris
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

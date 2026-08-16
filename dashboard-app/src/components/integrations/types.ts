export type IntegrationMode = 'managed' | 'custom_app';

export type ProviderId =
  | 'facebook'
  | 'tiktok'
  | 'google'
  | 'linkedin'
  | 'hubspot'
  | 'slack';

export interface ProviderMeta {
  id: ProviderId;
  name: string;
  category: 'social_messaging' | 'ads_crm' | 'calendar_workspace' | 'collaboration';
  icon: string;
  description: string;
  developerConsoleUrl: string;
  supportsCustomApp: boolean;
}

export interface TenantOAuthConfig {
  id: string;
  tenantId: string;
  provider: ProviderId;
  isCustomApp: boolean;
  customClientId?: string;
  customClientSecret?: string;
  customScopes?: string[];
  updatedAt: string;
}

export interface ConnectedAccount {
  id: string;
  tenantId: string;
  provider: ProviderId;
  accountName: string;
  accountAvatar?: string;
  email?: string;
  isCustomApp: boolean;
  status: 'active' | 'token_expired' | 'error';
  lastRefreshedAt: string;
}

export interface CustomAppFormValues {
  clientId: string;
  clientSecret: string;
  scopes?: string;
}

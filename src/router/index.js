import Vue from 'vue';
import VueRouter from 'vue-router';

// Lazy loaded views
import DashboardView from '@/views/DashboardView.vue';
import ProductListView from '@/views/products/ProductListView.vue';
import ProductDetailView from '@/views/products/ProductDetailView.vue';
import ProductFormView from '@/views/products/ProductFormView.vue';
import ComponentListView from '@/views/components/ComponentListView.vue';
import ComponentDetailView from '@/views/components/ComponentDetailView.vue';
import ComponentFormView from '@/views/components/ComponentFormView.vue';
import CategoryListView from '@/views/categories/CategoryListView.vue';
import ManufacturerListView from '@/views/manufacturers/ManufacturerListView.vue';
import SupplierListView from '@/views/suppliers/SupplierListView.vue';
import CustomerListView from '@/views/customers/CustomerListView.vue';
import UnitListView from '@/views/units/UnitListView.vue';
import UserListView from '@/views/users/UserListView.vue';
import RoleListView from '@/views/roles/RoleListView.vue';
import ProjectProgressView from '@/views/progress/ProjectProgressView.vue';
import ProductVersionListView from '@/views/versions/ProductVersionListView.vue';
import ProductVersionDetailView from '@/views/versions/ProductVersionDetailView.vue';
import ProductVersionFormView from '@/views/versions/ProductVersionFormView.vue';
import VersionCompareView from '@/views/versions/VersionCompareView.vue';
import HardwareRevisionListView from '@/views/revisions/HardwareRevisionListView.vue';
import HardwareRevisionDetailView from '@/views/revisions/HardwareRevisionDetailView.vue';
import EngineeringBOMWorkspaceView from '@/views/bom/EngineeringBOMWorkspaceView.vue';
import PrototypeDashboardView from '@/views/prototypes/PrototypeDashboardView.vue';
import PrototypeListView from '@/views/prototypes/PrototypeListView.vue';
import PrototypeCreateWizardView from '@/views/prototypes/PrototypeCreateWizardView.vue';
import PrototypeDetailView from '@/views/prototypes/PrototypeDetailView.vue';

// Phase 5: Testing & Validation Control
import TestingDashboardView from '@/views/testing/TestingDashboardView.vue';
import TestPlanRegistryView from '@/views/testing/TestPlanRegistryView.vue';
import TestPlanDetailView from '@/views/testing/TestPlanDetailView.vue';
import TestExecutionWorkspaceView from '@/views/testing/TestExecutionWorkspaceView.vue';
import EngineeringFindingsView from '@/views/testing/EngineeringFindingsView.vue';

// Phase 6: Engineering Change Management (ECR / ECO)
import ChangeDashboardView from '@/views/changes/ChangeDashboardView.vue';
import ECRRegistryView from '@/views/changes/ECRRegistryView.vue';
import ECRCreateWizardView from '@/views/changes/ECRCreateWizardView.vue';
import ECRDetailView from '@/views/changes/ECRDetailView.vue';
import ECORegistryView from '@/views/changes/ECORegistryView.vue';
import ECODetailView from '@/views/changes/ECODetailView.vue';
import TraceabilityMatrixView from '@/views/changes/TraceabilityMatrixView.vue';

// Cross-Cutting Platform Foundation: Data Import & Migration
import ImportDashboardView from '@/views/import/ImportDashboardView.vue';
import ImportWizardView from '@/views/import/ImportWizardView.vue';
import ImportReviewWorkspace from '@/views/import/ImportReviewWorkspace.vue';

// Developer & API Integrations
import ApiTokenView from '@/views/settings/ApiTokenView.vue';

Vue.use(VueRouter);

const routes = [
  {
    path: '/',
    name: 'Dashboard',
    component: DashboardView
  },
  // Cross-Cutting Platform Capability: Product Data Import & Migration
  {
    path: '/data-import',
    name: 'ImportDashboard',
    component: ImportDashboardView
  },
  {
    path: '/data-import/wizard/:id?',
    name: 'ImportWizard',
    component: ImportWizardView
  },
  {
    path: '/data-import/review/:id',
    name: 'ImportReview',
    component: ImportReviewWorkspace
  },
  // Project Management & Roadmap Progress
  {
    path: '/progress',
    name: 'ProjectProgress',
    component: ProjectProgressView
  },
  {
    path: '/project-progress',
    redirect: '/progress'
  },
  // Phase 2: Product Versions
  {
    path: '/versions',
    name: 'ProductVersions',
    component: ProductVersionListView
  },
  {
    path: '/versions/create',
    name: 'ProductVersionCreate',
    component: ProductVersionFormView
  },
  {
    path: '/versions/compare',
    name: 'VersionCompare',
    component: VersionCompareView
  },
  {
    path: '/versions/:id',
    name: 'ProductVersionDetail',
    component: ProductVersionDetailView
  },
  {
    path: '/versions/:id/edit',
    name: 'ProductVersionEdit',
    component: ProductVersionFormView
  },
  // Phase 2: Hardware Revisions
  {
    path: '/revisions',
    name: 'HardwareRevisions',
    component: HardwareRevisionListView
  },
  {
    path: '/revisions/:id',
    name: 'HardwareRevisionDetail',
    component: HardwareRevisionDetailView
  },
  // Phase 3: Engineering BOM Workspace
  {
    path: '/revisions/:id/bom',
    name: 'EngineeringBOMWorkspace',
    component: EngineeringBOMWorkspaceView
  },
  // Phase 4: Prototype & Engineering Build Management
  {
    path: '/prototypes',
    name: 'PrototypeDashboard',
    component: PrototypeDashboardView
  },
  {
    path: '/prototypes/list',
    name: 'PrototypeList',
    component: PrototypeListView
  },
  {
    path: '/prototypes/create',
    name: 'PrototypeCreate',
    component: PrototypeCreateWizardView
  },
  {
    path: '/prototypes/:id',
    name: 'PrototypeDetail',
    component: PrototypeDetailView
  },
  // Phase 5: Testing & Validation Control
  {
    path: '/testing',
    name: 'TestingDashboard',
    component: TestingDashboardView
  },
  {
    path: '/test-plans',
    name: 'TestPlanRegistry',
    component: TestPlanRegistryView
  },
  {
    path: '/test-plans/:id',
    name: 'TestPlanDetail',
    component: TestPlanDetailView
  },
  {
    path: '/test-sessions/:id',
    name: 'TestExecutionWorkspace',
    component: TestExecutionWorkspaceView
  },
  {
    path: '/findings',
    name: 'EngineeringFindings',
    component: EngineeringFindingsView
  },
  // Phase 6: Engineering Change Management (ECR / ECO)
  {
    path: '/changes',
    name: 'ChangeDashboard',
    component: ChangeDashboardView
  },
  {
    path: '/ecr',
    name: 'ECRRegistry',
    component: ECRRegistryView
  },
  {
    path: '/ecr/create',
    name: 'ECRCreate',
    component: ECRCreateWizardView
  },
  {
    path: '/ecr/:id',
    name: 'ECRDetail',
    component: ECRDetailView
  },
  {
    path: '/eco',
    name: 'ECORegistry',
    component: ECORegistryView
  },
  {
    path: '/eco/:id',
    name: 'ECODetail',
    component: ECODetailView
  },
  {
    path: '/changes/traceability',
    name: 'TraceabilityMatrix',
    component: TraceabilityMatrixView
  },
  // Products Module
  {
    path: '/products',
    name: 'Products',
    component: ProductListView
  },
  {
    path: '/products/create',
    name: 'ProductCreate',
    component: ProductFormView
  },
  {
    path: '/products/:id',
    name: 'ProductDetail',
    component: ProductDetailView
  },
  {
    path: '/products/:id/edit',
    name: 'ProductEdit',
    component: ProductFormView
  },
  // Components Module
  {
    path: '/components',
    name: 'Components',
    component: ComponentListView
  },
  {
    path: '/components/create',
    name: 'ComponentCreate',
    component: ComponentFormView
  },
  {
    path: '/components/:id',
    name: 'ComponentDetail',
    component: ComponentDetailView
  },
  {
    path: '/components/:id/edit',
    name: 'ComponentEdit',
    component: ComponentFormView
  },
  // Categories Module
  {
    path: '/categories',
    name: 'Categories',
    component: CategoryListView
  },
  // Manufacturers Module
  {
    path: '/manufacturers',
    name: 'Manufacturers',
    component: ManufacturerListView
  },
  // Suppliers Module
  {
    path: '/suppliers',
    name: 'Suppliers',
    component: SupplierListView
  },
  // Customers Module
  {
    path: '/customers',
    name: 'Customers',
    component: CustomerListView
  },
  // Units Module
  {
    path: '/units',
    name: 'Units',
    component: UnitListView
  },
  // Users Module
  {
    path: '/users',
    name: 'Users',
    component: UserListView
  },
  // Roles & Permissions Module
  {
    path: '/roles',
    name: 'Roles',
    component: RoleListView
  },
  // Developer & System Integrations: API Tokens
  {
    path: '/api-tokens',
    name: 'ApiTokens',
    component: ApiTokenView
  },
  // Authentication Module
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/LoginView.vue'),
    meta: { public: true }
  },
  // Catch-all
  {
    path: '*',
    redirect: '/'
  }
];

const router = new VueRouter({
  mode: 'hash', // Hash mode ensures flawless routing in dev server and standalone preview
  scrollBehavior() {
    return { x: 0, y: 0 };
  },
  routes
});

// Global Navigation Session Guard
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('iot_token') || sessionStorage.getItem('iot_token');
  const isPublic = to.matched.some(record => record.meta && record.meta.public);

  if (!isPublic && !token) {
    // Unauthenticated access to protected route -> redirect to login with destination preservation
    next({
      path: '/login',
      query: to.fullPath && to.fullPath !== '/' ? { redirect: to.fullPath } : undefined
    });
  } else if (to.path === '/login' && token) {
    // Authenticated user attempting to visit login page -> redirect to dashboard
    next({ path: '/' });
  } else {
    next();
  }
});

export default router;

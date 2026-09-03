<template>
  <div class="max-w-4xl mx-auto space-y-5 pb-16">
    <!-- Form Page Header -->
    <PageHeader
      :title="isEdit ? `Edit Product: ${form.code}` : 'Create IoT Product Master'"
      :subtitle="isEdit ? 'Update existing product master record' : 'Register a new enterprise IoT hardware product master'"
    >
      <template #actions>
        <button type="button" @click="cancel" class="btn btn-secondary">
          Cancel
        </button>
        <button type="button" @click="save('Draft')" class="btn btn-secondary">
          Save Draft
        </button>
        <button type="button" @click="save(form.status || 'Active')" class="btn btn-primary">
          <AppIcon name="check" size="xs" class="mr-1.5" />
          <span>{{ isEdit ? 'Save Changes' : 'Save Product' }}</span>
        </button>
      </template>
    </PageHeader>

    <form @submit.prevent="save(form.status || 'Active')" class="space-y-5">
      <!-- SECTION 00: PRODUCT TYPE SELECTION (Primary Architectural Choice) -->
      <div class="bg-white rounded-lg border border-slate-200 shadow-card p-5 space-y-3">
        <div class="flex items-center justify-between pb-3 border-b border-slate-200">
          <div class="flex items-center gap-2">
            <AppIcon name="layers" size="sm" class="text-brand-600" />
            <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">
              Product Classification Type *
            </h2>
          </div>
          <span class="text-3xs text-slate-500">Controls available modules, lifecycle rules, and composition capabilities</span>
        </div>

        <!-- 3 Product Type Cards -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
          <!-- CARD 1: MANUFACTURE PRODUCT -->
          <div
            @click="setProductType('MANUFACTURE')"
            :class="[
              'p-4 rounded-xl border-2 transition-all cursor-pointer flex flex-col justify-between select-none relative overflow-hidden',
              form.productType === 'MANUFACTURE' || form.productType === 'RND'
                ? 'border-emerald-500 bg-emerald-50/40 shadow-xs'
                : 'border-slate-200 bg-white hover:border-slate-300 hover:bg-slate-50/50'
            ]"
          >
            <div>
              <div class="flex items-center justify-between mb-2">
                <div class="w-8 h-8 rounded-lg bg-emerald-100/80 text-emerald-700 flex items-center justify-center font-bold">
                  <AppIcon name="cpu" size="xs" />
                </div>
                <span v-if="form.productType === 'MANUFACTURE' || form.productType === 'RND'" class="w-2.5 h-2.5 rounded-full bg-emerald-500"></span>
              </div>
              <h3 class="font-bold text-xs text-slate-900">MANUFACTURE PRODUCT</h3>
              <p class="text-3xs text-slate-500 mt-1 leading-relaxed">
                Internally manufactured engineering hardware with full lifecycle versioning.
              </p>
            </div>

            <div class="pt-3 mt-3 border-t border-slate-200/80 space-y-1 text-3xs">
              <div class="text-emerald-700 font-medium flex items-center gap-1">
                <span>✓</span> <span>Product Versions & Revisions</span>
              </div>
              <div class="text-emerald-700 font-medium flex items-center gap-1">
                <span>✓</span> <span>Future Engineering BOM</span>
              </div>
              <div class="text-slate-400 flex items-center gap-1">
                <span>✕</span> <span>No Resell Commercial Matrix</span>
              </div>
            </div>
          </div>

          <!-- CARD 2: TRADING PRODUCT -->
          <div
            @click="setProductType('TRADING')"
            :class="[
              'p-4 rounded-xl border-2 transition-all cursor-pointer flex flex-col justify-between select-none relative overflow-hidden',
              form.productType === 'TRADING'
                ? 'border-purple-500 bg-purple-50/40 shadow-xs'
                : 'border-slate-200 bg-white hover:border-slate-300 hover:bg-slate-50/50'
            ]"
          >
            <div>
              <div class="flex items-center justify-between mb-2">
                <div class="w-8 h-8 rounded-lg bg-purple-100/80 text-purple-700 flex items-center justify-center font-bold">
                  <AppIcon name="shopping-bag" size="xs" />
                </div>
                <span v-if="form.productType === 'TRADING'" class="w-2.5 h-2.5 rounded-full bg-purple-500"></span>
              </div>
              <h3 class="font-bold text-xs text-slate-900">TRADING PRODUCT</h3>
              <p class="text-3xs text-slate-500 mt-1 leading-relaxed">
                Externally sourced commercial hardware, third-party sensors, and ready-made devices.
              </p>
            </div>

            <div class="pt-3 mt-3 border-t border-slate-200/80 space-y-1 text-3xs">
              <div class="text-purple-700 font-medium flex items-center gap-1">
                <span>✓</span> <span>Commercial Pricing & Sourcing</span>
              </div>
              <div class="text-purple-700 font-medium flex items-center gap-1">
                <span>✓</span> <span>Manufacturer & Supplier Terms</span>
              </div>
              <div class="text-slate-400 flex items-center gap-1">
                <span>✕</span> <span>No Versioning / No BOM</span>
              </div>
            </div>
          </div>

          <!-- CARD 3: PROJECT PRODUCT -->
          <div
            @click="setProductType('PROJECT')"
            :class="[
              'p-4 rounded-xl border-2 transition-all cursor-pointer flex flex-col justify-between select-none relative overflow-hidden',
              form.productType === 'PROJECT'
                ? 'border-cyan-500 bg-cyan-50/40 shadow-xs'
                : 'border-slate-200 bg-white hover:border-slate-300 hover:bg-slate-50/50'
            ]"
          >
            <div>
              <div class="flex items-center justify-between mb-2">
                <div class="w-8 h-8 rounded-lg bg-cyan-100/80 text-cyan-700 flex items-center justify-center font-bold">
                  <AppIcon name="layers" size="xs" />
                </div>
                <span v-if="form.productType === 'PROJECT'" class="w-2.5 h-2.5 rounded-full bg-cyan-500"></span>
              </div>
              <h3 class="font-bold text-xs text-slate-900">PROJECT PRODUCT</h3>
              <p class="text-3xs text-slate-500 mt-1 leading-relaxed">
                Turnkey solution package composed of multiple products, sub-assemblies, and components.
              </p>
            </div>

            <div class="pt-3 mt-3 border-t border-slate-200/80 space-y-1 text-3xs">
              <div class="text-cyan-700 font-medium flex items-center gap-1">
                <span>✓</span> <span>System Architecture Tree</span>
              </div>
              <div class="text-cyan-700 font-medium flex items-center gap-1">
                <span>✓</span> <span>References Trading & Manufacture</span>
              </div>
              <div class="text-slate-400 flex items-center gap-1">
                <span>✕</span> <span>No Internal Manufacture BOM Duplication</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- SECTION 01: BASIC MASTER INFORMATION -->
      <div class="bg-white rounded-lg border border-slate-200 shadow-card p-5">
        <div class="flex items-center gap-2 pb-3 border-b border-slate-200 mb-4">
          <AppIcon name="product" size="sm" class="text-brand-600" />
          <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Master Data Information</h2>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- Product Code -->
          <div>
            <div class="flex items-center justify-between mb-1">
              <label class="form-label mb-0">Product Code *</label>
              <button
                type="button"
                @click="generateProductCode"
                class="text-3xs text-brand-600 hover:underline font-mono"
              >
                Auto-generate
              </button>
            </div>
            <input
              v-model="form.code"
              type="text"
              required
              placeholder="e.g. PRD-2024-009"
              class="form-input font-mono uppercase"
            />
          </div>

          <!-- Product Name -->
          <div>
            <label class="form-label">Product Name *</label>
            <input
              v-model="form.name"
              type="text"
              required
              placeholder="e.g. Smart Edge Industrial Gateway (GW-400X)"
              class="form-input font-medium"
            />
          </div>

          <!-- Product Category -->
          <div>
            <label class="form-label">Product Category *</label>
            <select v-model="form.category" required class="form-select">
              <option value="" disabled>Select Category</option>
              <option value="Edge Gateways">Edge Gateways</option>
              <option value="Sensors & Nodes">Sensors & Nodes</option>
              <option value="Power Systems">Power Systems</option>
              <option value="Asset Tracking">Asset Tracking</option>
              <option value="Industrial Controllers">Industrial Controllers</option>
              <option value="Integrated Solutions">Integrated Solutions</option>
              <option value="Wearable & Telemetry">Wearable & Telemetry</option>
            </select>
          </div>

          <!-- Status -->
          <div>
            <label class="form-label">Lifecycle Status *</label>
            <select v-model="form.status" required class="form-select">
              <option value="Draft">Draft</option>
              <option value="Active">Active</option>
              <option value="In Review">In Review</option>
              <option value="Obsolete">Obsolete</option>
            </select>
          </div>

          <!-- Project Lead -->
          <div>
            <label class="form-label">Lead Architect / Owner</label>
            <select v-model="form.projectLead" class="form-select">
              <option value="">Select Project Lead</option>
              <option v-for="u in store.users" :key="u.id" :value="u.name">
                {{ u.name }} ({{ getRoleName(u.role) }})
              </option>
            </select>
          </div>

          <!-- Target Market -->
          <div>
            <label class="form-label">Target Market / Application</label>
            <input
              v-model="form.targetMarket"
              type="text"
              placeholder="e.g. Smart City & Factory Automation"
              class="form-input"
            />
          </div>
        </div>
      </div>

      <!-- SECTION: PRODUCT PHOTOS & GALLERY (MAX 5 IMAGES) -->
      <div class="bg-white rounded-lg border border-slate-200 shadow-card p-5 space-y-4">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 pb-3 border-b border-slate-100">
          <div class="flex items-center gap-2">
            <div class="w-8 h-8 rounded-lg bg-brand-50 text-brand-600 flex items-center justify-center">
              <AppIcon name="camera" size="xs" />
            </div>
            <div>
              <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Product Photos & Gallery</h2>
              <p class="text-3xs text-slate-500">Upload high-res product photos, PCB renders, enclosures, or technical drawings (Max 5 photos).</p>
            </div>
          </div>
          <span
            class="text-3xs font-mono font-bold px-2.5 py-0.5 rounded-full border self-start sm:self-auto"
            :class="productImages.length >= 5 ? 'bg-amber-50 text-amber-700 border-amber-200' : 'bg-brand-50 text-brand-700 border-brand-200'"
          >
            {{ productImages.length }} / 5 Photos
          </span>
        </div>

        <!-- UPLOAD DROPZONE -->
        <div
          v-if="productImages.length < 5"
          @click="triggerFileInput"
          @dragover.prevent="isDragging = true"
          @dragleave.prevent="isDragging = false"
          @drop.prevent="handleFileDrop"
          class="border-2 border-dashed rounded-xl p-6 text-center cursor-pointer transition-all duration-200"
          :class="isDragging ? 'border-brand-500 bg-brand-50/50 scale-[0.99]' : 'border-slate-300 hover:border-brand-400 bg-slate-50/60 hover:bg-white'"
        >
          <input
            ref="fileInputRef"
            type="file"
            multiple
            accept="image/png,image/jpeg,image/webp,image/svg+xml"
            class="hidden"
            @change="handleFileInputChange"
          />

          <div v-if="isUploading" class="flex flex-col items-center justify-center gap-2 py-2">
            <AppIcon name="refresh" size="md" class="animate-spin text-brand-600" />
            <span class="text-xs font-bold text-slate-700">Uploading photo(s)...</span>
          </div>

          <div v-else class="flex flex-col items-center justify-center gap-2">
            <div class="w-10 h-10 rounded-full bg-brand-100 text-brand-600 flex items-center justify-center shadow-xs">
              <AppIcon name="upload" size="sm" />
            </div>
            <div>
              <span class="text-xs font-bold text-slate-800">Click to browse</span>
              <span class="text-xs text-slate-500"> or drag & drop images here</span>
            </div>
            <p class="text-3xs text-slate-400">
              Supports PNG, JPG, WebP up to 10MB each ({{ 5 - productImages.length }} slots remaining)
            </p>
          </div>
        </div>

        <!-- GALLERY PREVIEW GRID -->
        <div v-if="productImages.length > 0" class="space-y-3">
          <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-3 pt-1">
            <div
              v-for="(img, idx) in productImages"
              :key="idx"
              class="group relative rounded-xl border overflow-hidden bg-slate-100 aspect-square shadow-2xs transition-all hover:shadow-md"
              :class="idx === 0 ? 'border-brand-500 ring-2 ring-brand-500/30' : 'border-slate-200'"
            >
              <img
                :src="getImageUrl(img)"
                alt="Product photo"
                class="w-full h-full object-cover"
                @error="onImgPreviewError(idx)"
              />

              <!-- Primary Cover Badge -->
              <div
                v-if="idx === 0"
                class="absolute top-1.5 left-1.5 bg-brand-600 text-white text-3xs font-bold px-1.5 py-0.5 rounded shadow-xs flex items-center gap-1"
              >
                <AppIcon name="check" size="2xs" />
                <span>Cover Photo</span>
              </div>

              <!-- Slot number badge for other photos -->
              <div
                v-else
                class="absolute top-1.5 left-1.5 bg-slate-900/70 text-white text-3xs font-mono font-bold px-1.5 py-0.5 rounded backdrop-blur-xs"
              >
                #{{ idx + 1 }}
              </div>

              <!-- Overlay Actions on Hover -->
              <div class="absolute inset-0 bg-slate-900/60 opacity-0 group-hover:opacity-100 transition-opacity flex flex-col items-center justify-center gap-2 p-2">
                <button
                  v-if="idx !== 0"
                  type="button"
                  @click.stop="setAsPrimaryImage(idx)"
                  class="btn btn-secondary py-1 px-2 text-3xs font-bold bg-white text-slate-900 hover:bg-slate-100 w-full"
                  title="Set as main catalog cover"
                >
                  Set as Cover
                </button>

                <button
                  type="button"
                  @click.stop="removeImage(idx)"
                  class="btn btn-danger py-1 px-2 text-3xs font-bold w-full flex items-center justify-center gap-1"
                  title="Remove photo"
                >
                  <AppIcon name="trash" size="2xs" />
                  <span>Delete</span>
                </button>
              </div>
            </div>
          </div>
        </div>

        <div v-else class="text-center py-3 bg-slate-50/50 rounded-lg border border-slate-100">
          <p class="text-3xs text-slate-400">No photos uploaded yet. The first photo uploaded will automatically be used as the catalog cover.</p>
        </div>

        <!-- Direct URL Option -->
        <div class="pt-2 border-t border-slate-100">
          <details class="group text-xs">
            <summary class="cursor-pointer text-3xs font-semibold text-slate-500 hover:text-slate-800 list-none flex items-center gap-1.5">
              <AppIcon name="link" size="2xs" class="text-slate-400" />
              <span>Or add image via direct Web URL</span>
            </summary>
            <div class="flex items-center gap-2 mt-2">
              <input
                v-model="directUrlInput"
                type="url"
                placeholder="https://example.com/product-image.jpg"
                class="form-input text-xs flex-1"
                @keyup.enter="addDirectUrlImage"
              />
              <button
                type="button"
                @click="addDirectUrlImage"
                class="btn btn-secondary text-xs shrink-0"
                :disabled="productImages.length >= 5 || !directUrlInput"
              >
                Add URL
              </button>
            </div>
          </details>
        </div>
      </div>

      <!-- SECTION 02: DYNAMIC SPECIFICATIONS ACCORDING TO PRODUCT TYPE -->

      <!-- A: IF TRADING PRODUCT -> COMMERCIAL SOURCING & MULTI-CURRENCY PRICING -->
      <div v-if="form.productType === 'TRADING'" class="bg-white rounded-lg border border-purple-200 shadow-card p-5 space-y-6">
        <div class="flex items-center justify-between pb-3 border-b border-purple-100">
          <div class="flex items-center gap-2">
            <AppIcon name="shopping-bag" size="sm" class="text-purple-600" />
            <h2 class="text-xs font-bold uppercase tracking-wider text-purple-900">Trading Commercial & Sourcing Parameters</h2>
          </div>
          <span class="text-3xs font-mono font-bold px-2.5 py-0.5 rounded-full bg-purple-50 text-purple-700 border border-purple-200">
            MULTI-CURRENCY PRICING ENGINE
          </span>
        </div>

        <!-- Sourcing Details -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-xs">
          <div>
            <label class="form-label">Manufacturer Name</label>
            <input v-model="tradingForm.manufacturerName" type="text" class="form-input" placeholder="e.g. Vaisala Instruments Oyj" />
          </div>
          <div>
            <label class="form-label">Manufacturer Part Number (MPN)</label>
            <input v-model="tradingForm.manufacturerPartNumber" type="text" class="form-input font-mono" placeholder="e.g. WXT536-A" />
          </div>
          <div>
            <label class="form-label">Authorized Supplier / Distributor</label>
            <input v-model="tradingForm.supplierName" type="text" class="form-input" placeholder="e.g. DigiKey Electronics" />
          </div>
          <div>
            <label class="form-label">Supplier SKU / Order Number</label>
            <input v-model="tradingForm.supplierPartNumber" type="text" class="form-input font-mono" placeholder="e.g. DK-WXT536-ND" />
          </div>
        </div>

        <!-- MULTI-CURRENCY PRICING & PROFIT ENGINE WORKSPACE -->
        <div class="p-4.5 rounded-xl bg-purple-50/40 border border-purple-200/80 space-y-5">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <AppIcon name="dollar-sign" size="xs" class="text-purple-700" />
              <h3 class="text-xs font-bold uppercase tracking-wider text-purple-950">Commercial Pricing & Margin Model</h3>
            </div>
            <div class="text-3xs text-purple-600 font-medium">Selling price auto-derived from cost + profit margin</div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <!-- 1. PURCHASE PRICE -->
            <div class="p-3.5 rounded-lg bg-white border border-purple-200/80 shadow-2xs space-y-3">
              <div class="flex items-center justify-between">
                <span class="text-3xs uppercase font-bold text-slate-500">1. Purchase Price (Cost)</span>
                <span class="text-3xs px-1.5 py-0.2 rounded font-mono font-bold bg-slate-100 text-slate-700">{{ tradingForm.purchaseCurrency }}</span>
              </div>
              <div class="grid grid-cols-5 gap-2">
                <input
                  v-model.number="tradingForm.purchasePrice"
                  @input="recalculatePricing"
                  type="number"
                  step="0.01"
                  min="0"
                  class="form-input col-span-3 font-mono font-bold text-xs"
                  placeholder="0.00"
                />
                <select
                  v-model="tradingForm.purchaseCurrency"
                  @change="onPurchaseCurrencyChange"
                  class="form-select col-span-2 font-mono text-xs"
                >
                  <option value="USD">USD ($)</option>
                  <option value="IDR">IDR (Rp)</option>
                  <option value="EUR">EUR (€)</option>
                  <option value="SGD">SGD ($)</option>
                </select>
              </div>
              <div class="text-3xs text-slate-400 font-mono flex justify-between">
                <span>Base Unit Cost:</span>
                <span class="font-bold text-slate-700">{{ formatCurrency(tradingForm.purchasePrice, tradingForm.purchaseCurrency) }}</span>
              </div>
            </div>

            <!-- 2. PROFIT MARGIN -->
            <div class="p-3.5 rounded-lg bg-white border border-purple-200/80 shadow-2xs space-y-3">
              <div class="flex items-center justify-between">
                <span class="text-3xs uppercase font-bold text-slate-500">2. Profit Margin</span>
                <span class="text-3xs font-mono font-bold text-brand-700 px-1.5 py-0.2 rounded bg-brand-50 border border-brand-200">+{{ tradingForm.profitMarginPercentage }}%</span>
              </div>
              <div class="relative">
                <input
                  v-model.number="tradingForm.profitMarginPercentage"
                  @input="recalculatePricing"
                  type="number"
                  step="1"
                  min="0"
                  max="500"
                  class="form-input font-mono font-bold text-xs pr-7"
                  placeholder="25"
                />
                <span class="absolute right-3 top-2.5 text-xs text-slate-400 font-bold">%</span>
              </div>
              <div class="text-3xs text-slate-400 leading-tight">
                Formula: <span class="font-mono text-slate-600">Cost × (1 + {{ tradingForm.profitMarginPercentage }}%)</span>
              </div>
            </div>

            <!-- 3. SELLING PRICE -->
            <div class="p-3.5 rounded-lg bg-purple-100/50 border border-purple-300 shadow-2xs space-y-3">
              <div class="flex items-center justify-between">
                <span class="text-3xs uppercase font-bold text-purple-900">3. Calculated Selling Price</span>
                <span class="text-3xs px-1.5 py-0.2 rounded font-mono font-bold bg-purple-200 text-purple-900">{{ tradingForm.sellingCurrency }}</span>
              </div>
              <div class="grid grid-cols-5 gap-2">
                <input
                  :value="tradingForm.sellingPrice"
                  readonly
                  type="number"
                  class="form-input col-span-3 font-mono font-bold text-xs bg-white text-purple-900 border-purple-300 cursor-not-allowed shadow-inner"
                />
                <select
                  v-model="tradingForm.sellingCurrency"
                  @change="recalculatePricing"
                  class="form-select col-span-2 font-mono text-xs border-purple-300 bg-white"
                >
                  <option value="USD">USD ($)</option>
                  <option value="IDR">IDR (Rp)</option>
                  <option value="EUR">EUR (€)</option>
                  <option value="SGD">SGD ($)</option>
                </select>
              </div>
              <div class="text-3xs text-purple-900 font-mono flex justify-between">
                <span>Final Catalog MSRP:</span>
                <span class="font-bold text-purple-950">{{ formatCurrency(tradingForm.sellingPrice, tradingForm.sellingCurrency) }}</span>
              </div>
            </div>
          </div>

          <!-- EXCHANGE RATE & IDR EQUIVALENT REFERENCE CONTROLS -->
          <div class="p-4.5 rounded-xl bg-white border border-purple-200 shadow-2xs space-y-4">
            <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 pb-3 border-b border-slate-100">
              <div class="flex items-center gap-2">
                <AppIcon name="refresh" size="xs" class="text-indigo-600" />
                <span class="text-xs font-bold text-slate-800 uppercase tracking-wider">Exchange Rate (Kurs Dollar / IDR)</span>
              </div>

              <!-- Live Rate Status & Mode Badge -->
              <div class="flex items-center gap-2">
                <span
                  class="text-3xs font-mono px-2.5 py-0.5 rounded-full font-bold border"
                  :class="tradingForm.exchangeRateMode === 'MANUAL' ? 'bg-amber-50 text-amber-800 border-amber-200' : 'bg-emerald-50 text-emerald-800 border-emerald-200'"
                >
                  {{ tradingForm.exchangeRateMode === 'MANUAL' ? 'MANUAL / CUSTOM RATE' : 'LIVE MARKET RATE' }}
                </span>
                <button
                  type="button"
                  @click="refreshExchangeRate"
                  class="btn btn-secondary py-1 px-2.5 text-3xs font-semibold flex items-center gap-1.5 hover:bg-slate-100"
                  title="Fetch latest live market rate from ExchangeRate-API"
                >
                  <AppIcon name="refresh" size="2xs" class="text-brand-600" />
                  <span>Update Live Rate</span>
                </button>
              </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-3 gap-4 items-stretch">
              <!-- DIRECT MANUAL / LIVE INPUT FIELD -->
              <div class="p-3.5 rounded-lg bg-indigo-50/50 border border-indigo-200 space-y-2.5">
                <div class="flex items-center justify-between">
                  <label class="text-3xs uppercase font-bold text-slate-700">Effective Kurs (1 USD = Rp)</label>
                  <span class="text-3xs font-mono font-bold text-indigo-700 truncate max-w-[120px]" :title="tradingForm.exchangeRateProvider">
                    {{ tradingForm.exchangeRateProvider || 'ExchangeRate-API' }}
                  </span>
                </div>

                <div class="relative">
                  <span class="absolute left-3 top-2.5 text-xs font-mono font-bold text-slate-400">Rp</span>
                  <input
                    v-model.number="tradingForm.manualExchangeRate"
                    @input="onManualRateInput"
                    type="number"
                    step="1"
                    min="1"
                    class="form-input pl-9 font-mono font-bold text-xs bg-white text-slate-900 border-indigo-300 focus:border-brand-500 focus:ring-brand-500 shadow-inner"
                    placeholder="17749"
                  />
                </div>

                <!-- Quick Action Button to Sync Live Rate -->
                <div class="flex items-center justify-between text-3xs pt-1.5 border-t border-indigo-100">
                  <button
                    type="button"
                    @click="useLiveRate"
                    class="text-brand-700 hover:text-brand-900 font-bold flex items-center gap-1"
                    title="Click to apply live market rate"
                  >
                    <AppIcon name="check" size="2xs" class="text-emerald-600" />
                    <span>Sync Live: Rp {{ Number(tradingForm.liveExchangeRate || 17749).toLocaleString('id-ID') }}</span>
                  </button>
                  <span class="text-slate-400 font-mono text-3xs">{{ tradingForm.exchangeRateMode === 'MANUAL' ? 'Locked' : 'Auto' }}</span>
                </div>
              </div>

              <!-- IDR Purchase Reference Card -->
              <div class="p-3.5 rounded-lg bg-slate-50 border border-slate-200 text-xs flex flex-col justify-between">
                <div>
                  <span class="text-3xs uppercase font-bold text-slate-500 block">Purchase Cost in IDR</span>
                  <div class="font-mono font-bold text-base text-slate-900 mt-1">
                    Rp {{ Number(tradingForm.purchasePriceInIDR).toLocaleString('id-ID') }}
                  </div>
                </div>
                <span class="text-3xs text-slate-400 mt-2 block font-mono">
                  Base: {{ formatCurrency(tradingForm.purchasePrice, tradingForm.purchaseCurrency) }}
                </span>
              </div>

              <!-- IDR Selling Reference Card -->
              <div class="p-3.5 rounded-lg bg-emerald-50/80 border border-emerald-200 text-xs flex flex-col justify-between">
                <div>
                  <span class="text-3xs uppercase font-bold text-emerald-800 block">Selling Price in IDR</span>
                  <div class="font-mono font-bold text-base text-emerald-900 mt-1">
                    Rp {{ Number(tradingForm.sellingPriceInIDR).toLocaleString('id-ID') }}
                  </div>
                </div>
                <span class="text-3xs text-emerald-700 mt-2 block font-mono">
                  MSRP: {{ formatCurrency(tradingForm.sellingPrice, tradingForm.sellingCurrency) }} (+{{ tradingForm.profitMarginPercentage }}%)
                </span>
              </div>
            </div>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-xs">
          <div>
            <label class="form-label">Procurement Lead Time (Days)</label>
            <input v-model.number="tradingForm.leadTimeDays" type="number" class="form-input font-mono" placeholder="14" />
          </div>
          <div>
            <label class="form-label">Warranty Information</label>
            <input v-model="tradingForm.warrantyInformation" type="text" class="form-input" placeholder="e.g. 24 Months Factory Limited Warranty" />
          </div>
        </div>
      </div>

      <!-- B: IF MANUFACTURE PRODUCT -> ELECTRICAL & MECHANICAL SPECIFICATIONS -->
      <div v-else-if="form.productType === 'MANUFACTURE' || form.productType === 'RND'" class="bg-white rounded-lg border border-slate-200 shadow-card p-5">
        <div class="flex items-center gap-2 pb-3 border-b border-slate-200 mb-4">
          <AppIcon name="cpu" size="sm" class="text-brand-600" />
          <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Hardware & Environmental Baseline</h2>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label class="form-label">Power Rating / Supply</label>
            <input
              v-model="form.powerRating"
              type="text"
              placeholder="e.g. 9-36V DC, 15W Max"
              class="form-input font-mono"
            />
          </div>

          <div>
            <label class="form-label">Ingress Protection (IP)</label>
            <input
              v-model="form.ingressProtection"
              type="text"
              placeholder="e.g. IP67 / NEMA 4X"
              class="form-input font-mono"
            />
          </div>

          <div>
            <label class="form-label">Operating Temperature</label>
            <input
              v-model="form.operatingTemp"
              type="text"
              placeholder="e.g. -40°C to +85°C"
              class="form-input font-mono"
            />
          </div>
        </div>
      </div>

      <!-- SECTION 03: DESCRIPTION -->
      <div class="bg-white rounded-lg border border-slate-200 shadow-card p-5">
        <div class="flex items-center gap-2 pb-3 border-b border-slate-200 mb-4">
          <AppIcon name="file-text" size="sm" class="text-brand-600" />
          <h2 class="text-xs font-bold uppercase tracking-wider text-slate-800">Description & Scope</h2>
        </div>

        <div class="space-y-4">
          <div>
            <label class="form-label">Short Summary</label>
            <input
              v-model="form.description"
              type="text"
              placeholder="Brief 1-sentence product summary"
              class="form-input"
            />
          </div>

          <div>
            <label class="form-label">Detailed Scope Description</label>
            <textarea
              v-model="form.longDescription"
              rows="4"
              placeholder="Detailed specifications, primary interfaces, operating constraints, and system notes..."
              class="form-textarea"
            ></textarea>
          </div>
        </div>
      </div>
    </form>
  </div>
</template>

<script>
import { useMasterStore } from '@/stores/masterStore';
import { exchangeRateApi } from '@/api/exchangeRate';
import { uploadApi } from '@/api/upload';
import { getImageUrl } from '@/utils/image';
import PageHeader from '@/components/common/PageHeader.vue';
import AppIcon from '@/components/common/AppIcon.vue';

export default {
  name: 'ProductFormView',
  components: {
    PageHeader,
    AppIcon
  },
  data() {
    return {
      form: {
        id: null,
        code: '',
        name: '',
        productType: 'MANUFACTURE',
        category: 'Sensors & Nodes',
        status: 'Active',
        projectLead: '',
        targetMarket: '',
        description: '',
        longDescription: '',
        powerRating: '',
        ingressProtection: '',
        operatingTemp: '',
        image: '',
        images: '',
        certifications: '["CE","FCC","RoHS"]'
      },
      productImages: [],
      isUploading: false,
      isDragging: false,
      directUrlInput: '',
      tradingForm: {
        manufacturerName: '',
        manufacturerPartNumber: '',
        supplierName: '',
        supplierPartNumber: '',
        purchasePrice: 100,
        purchaseCurrency: 'USD',
        profitMarginPercentage: 25,
        sellingPrice: 125,
        sellingCurrency: 'USD',
        exchangeRateMode: 'AUTO',
        liveExchangeRate: 17749,
        manualExchangeRate: 17749,
        effectiveExchangeRate: 17749,
        exchangeRateProvider: 'ExchangeRate-API (Live)',
        purchasePriceInIDR: 1774900,
        sellingPriceInIDR: 2218625,
        currency: 'USD',
        leadTimeDays: 14,
        warrantyInformation: '24 Months Factory Limited Warranty',
        notes: ''
      }
    };
  },
  computed: {
    store() {
      return useMasterStore();
    },
    isEdit() {
      return Boolean(this.$route.params.id);
    },
    productId() {
      return this.$route.params.id;
    }
  },
  async mounted() {
    if (this.store.users.length === 0) {
      await this.store.fetchUsers();
    }

    // Load initial live exchange rate
    await this.fetchInitialExchangeRate();

    if (this.isEdit) {
      const p = this.store.products.find(item => item.id === this.productId || item.code === this.productId);
      if (p) {
        this.form = { ...p };

        // Parse images array
        if (p.images) {
          try {
            const parsed = typeof p.images === 'string' ? JSON.parse(p.images) : p.images;
            if (Array.isArray(parsed)) {
              this.productImages = parsed.filter(Boolean);
            }
          } catch (e) {
            console.warn('Could not parse product images JSON:', e);
          }
        }
        if (this.productImages.length === 0 && p.image) {
          this.productImages = [p.image];
        }

        if (this.form.productType === 'TRADING') {
          const detail = await this.store.fetchTradingDetail(p.id);
          if (detail) {
            this.tradingForm = {
              ...this.tradingForm,
              ...detail,
              purchaseCurrency: detail.purchaseCurrency || detail.currency || 'USD',
              sellingCurrency: detail.sellingCurrency || detail.purchaseCurrency || detail.currency || 'USD',
              profitMarginPercentage: detail.profitMarginPercentage || 25,
              exchangeRateMode: detail.exchangeRateMode || 'AUTO',
              manualExchangeRate: detail.manualExchangeRate || detail.effectiveExchangeRate || 17749
            };
            await this.recalculatePricing();
          }
        }
      }
    } else {
      await this.recalculatePricing();
    }
  },
  methods: {
    getImageUrl(url) {
      return getImageUrl(url);
    },
    formatCurrency(amount, currency) {
      if (amount === null || amount === undefined) return '—';
      const curr = (currency || 'USD').toUpperCase();
      if (curr === 'IDR') {
        return new Intl.NumberFormat('id-ID', {
          style: 'currency',
          currency: 'IDR',
          minimumFractionDigits: 0,
          maximumFractionDigits: 0
        }).format(amount);
      }
      return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: curr,
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
      }).format(amount);
    },
    async fetchInitialExchangeRate() {
      try {
        const res = await exchangeRateApi.getLatestRate('USD', 'IDR');
        if (res && res.data) {
          this.tradingForm.liveExchangeRate = res.data.rate || 16500;
          this.tradingForm.exchangeRateProvider = res.data.provider || 'Frankfurter API';
          if (!this.tradingForm.manualExchangeRate) {
            this.tradingForm.manualExchangeRate = this.tradingForm.liveExchangeRate;
          }
        }
      } catch (err) {
        console.warn('Could not fetch initial exchange rate, using fallback:', err);
      }
    },
    async refreshExchangeRate() {
      try {
        const res = await exchangeRateApi.refreshRate('USD', 'IDR');
        if (res && res.data) {
          this.tradingForm.liveExchangeRate = res.data.rate;
          this.tradingForm.exchangeRateProvider = res.data.provider;
          await this.recalculatePricing();
          this.store.showToast('Exchange Rate Updated', `Live USD/IDR rate updated to Rp ${Number(res.data.rate).toLocaleString('id-ID')}`, 'success');
        }
      } catch (err) {
        this.store.showToast('Rate Refresh Failed', 'Unable to refresh live rate from external provider.', 'warning');
      }
    },
    onPurchaseCurrencyChange() {
      // Default selling currency to match purchase currency
      this.tradingForm.sellingCurrency = this.tradingForm.purchaseCurrency;
      this.recalculatePricing();
    },
    onManualRateInput() {
      this.tradingForm.exchangeRateMode = 'MANUAL';
      const rate = Number(this.tradingForm.manualExchangeRate) || 0;
      this.tradingForm.effectiveExchangeRate = rate > 0 ? rate : (this.tradingForm.liveExchangeRate || 17749);
      this.recalculatePricing();
    },
    useLiveRate() {
      this.tradingForm.exchangeRateMode = 'AUTO';
      const live = this.tradingForm.liveExchangeRate || 17749;
      this.tradingForm.manualExchangeRate = live;
      this.tradingForm.effectiveExchangeRate = live;
      this.recalculatePricing();
      this.store.showToast('Live Rate Synced', `Exchange rate updated to live market Rp ${Number(live).toLocaleString('id-ID')}`, 'success');
    },
    async recalculatePricing() {
      try {
        const res = await exchangeRateApi.calculatePricing({
          purchasePrice: Number(this.tradingForm.purchasePrice) || 0,
          purchaseCurrency: this.tradingForm.purchaseCurrency || 'USD',
          profitMarginPercentage: Number(this.tradingForm.profitMarginPercentage) || 0,
          sellingCurrency: this.tradingForm.sellingCurrency || 'USD',
          exchangeRateMode: this.tradingForm.exchangeRateMode || 'AUTO',
          manualExchangeRate: Number(this.tradingForm.manualExchangeRate) || 0
        });

        if (res && res.data) {
          this.tradingForm.sellingPrice = res.data.sellingPrice;
          this.tradingForm.liveExchangeRate = res.data.liveExchangeRate;
          this.tradingForm.effectiveExchangeRate = res.data.effectiveExchangeRate;
          this.tradingForm.exchangeRateProvider = res.data.exchangeRateProvider;
          this.tradingForm.purchasePriceInIDR = res.data.purchasePriceInIDR;
          this.tradingForm.sellingPriceInIDR = res.data.sellingPriceInIDR;
        }
      } catch (err) {
        // Fallback local math
        const cost = Number(this.tradingForm.purchasePrice) || 0;
        const margin = Number(this.tradingForm.profitMarginPercentage) || 0;
        const baseSelling = cost * (1 + margin / 100);
        const effRate = this.tradingForm.exchangeRateMode === 'MANUAL' && this.tradingForm.manualExchangeRate > 0
          ? this.tradingForm.manualExchangeRate
          : (this.tradingForm.liveExchangeRate || 16500);

        this.tradingForm.effectiveExchangeRate = effRate;

        if (this.tradingForm.purchaseCurrency === this.tradingForm.sellingCurrency) {
          this.tradingForm.sellingPrice = Math.round(baseSelling * 100) / 100;
        } else if (this.tradingForm.purchaseCurrency === 'USD' && this.tradingForm.sellingCurrency === 'IDR') {
          this.tradingForm.sellingPrice = Math.round(baseSelling * effRate);
        } else if (this.tradingForm.purchaseCurrency === 'IDR' && this.tradingForm.sellingCurrency === 'USD') {
          this.tradingForm.sellingPrice = Math.round((baseSelling / effRate) * 100) / 100;
        } else {
          this.tradingForm.sellingPrice = Math.round(baseSelling * 100) / 100;
        }

        this.tradingForm.purchasePriceInIDR = this.tradingForm.purchaseCurrency === 'IDR' ? cost : Math.round(cost * effRate);
        this.tradingForm.sellingPriceInIDR = this.tradingForm.sellingCurrency === 'IDR' ? this.tradingForm.sellingPrice : Math.round(this.tradingForm.sellingPrice * effRate);
      }
    },
    triggerFileInput() {
      if (this.$refs.fileInputRef) {
        this.$refs.fileInputRef.click();
      }
    },
    handleFileInputChange(e) {
      const files = Array.from(e.target.files || []);
      if (files.length > 0) {
        this.uploadFiles(files);
      }
      e.target.value = ''; // reset
    },
    handleFileDrop(e) {
      this.isDragging = false;
      const files = Array.from(e.dataTransfer.files || []).filter(f => f.type.startsWith('image/'));
      if (files.length > 0) {
        this.uploadFiles(files);
      }
    },
    async uploadFiles(files) {
      const remainingSlots = 5 - this.productImages.length;
      if (remainingSlots <= 0) {
        this.store.showToast('Maximum Photos Reached', 'You can attach up to 5 photos per product.', 'warning');
        return;
      }

      const filesToUpload = files.slice(0, remainingSlots);
      if (files.length > remainingSlots) {
        this.store.showToast('Photo Limit', `Only ${remainingSlots} photo(s) could be added (max 5 photos).`, 'info');
      }

      this.isUploading = true;
      let successCount = 0;

      for (const file of filesToUpload) {
        try {
          const res = await uploadApi.uploadFile(file);
          const uploadedUrl = res.data?.data?.url || res.data?.url || res.data?.filePath;
          if (uploadedUrl) {
            this.productImages.push(uploadedUrl);
            successCount++;
          }
        } catch (err) {
          console.error('Failed to upload product photo:', err);
          this.store.showToast('Upload Failed', `Failed to upload ${file.name}: ${err.message}`, 'danger');
        }
      }

      this.isUploading = false;
      if (successCount > 0) {
        this.store.showToast('Photos Uploaded', `${successCount} photo(s) added to product gallery.`, 'success');
      }
    },
    setAsPrimaryImage(idx) {
      if (idx > 0 && idx < this.productImages.length) {
        const selected = this.productImages.splice(idx, 1)[0];
        this.productImages.unshift(selected);
        this.store.showToast('Cover Photo Set', 'Primary catalog cover photo updated.', 'info');
      }
    },
    removeImage(idx) {
      if (idx >= 0 && idx < this.productImages.length) {
        this.productImages.splice(idx, 1);
        this.store.showToast('Photo Removed', 'Photo removed from gallery.', 'info');
      }
    },
    addDirectUrlImage() {
      const url = (this.directUrlInput || '').trim();
      if (!url) return;
      if (this.productImages.length >= 5) {
        this.store.showToast('Maximum Photos Reached', 'You can attach up to 5 photos per product.', 'warning');
        return;
      }
      this.productImages.push(url);
      this.directUrlInput = '';
      this.store.showToast('Photo Added', 'Image URL added to gallery.', 'success');
    },
    onImgPreviewError(idx) {
      // Fallback placeholder if broken
      console.warn(`Product photo #${idx + 1} failed to load.`);
    },
    getRoleName(role) {
      if (!role) return '';
      if (typeof role === 'object') return role.name || '';
      return role;
    },
    setProductType(type) {
      this.form.productType = type;
      if (!this.isEdit && !this.form.code) {
        this.generateProductCode();
      }
    },
    generateProductCode() {
      const count = this.store.products.length + 1;
      let prefix = 'PRD';
      if (this.form.productType === 'TRADING') {
        prefix = 'PRD-TRD';
      } else if (this.form.productType === 'PROJECT') {
        prefix = 'PRD-PRJ';
      }
      this.form.code = `${prefix}-2024-${String(count).padStart(3, '0')}`;
    },
    cancel() {
      if (this.isEdit) {
        this.$router.push(`/products/${this.productId}`);
      } else {
        this.$router.push('/products');
      }
    },
    async save(statusOverride) {
      if (!this.form.code || !this.form.name || !this.form.category) {
        alert('Please fill all required fields: Code, Name, Category.');
        return;
      }

      this.form.status = statusOverride || this.form.status;

      // Sync photos into form
      this.form.images = JSON.stringify(this.productImages);
      this.form.image = this.productImages.length > 0 ? this.productImages[0] : '';

      const payload = {
        ...this.form,
        images: JSON.stringify(this.productImages),
        image: this.productImages.length > 0 ? this.productImages[0] : '',
        tradingDetail: this.form.productType === 'TRADING' ? this.tradingForm : null
      };

      try {
        const saved = await this.store.saveProduct(payload);
        const targetId = saved.id || this.form.id || this.form.code;

        if (this.form.productType === 'TRADING' && this.isEdit) {
          await this.store.saveTradingDetail(targetId, this.tradingForm);
        }

        this.store.showToast('Product Saved', `Product '${this.form.name}' saved successfully.`, 'success');
        this.$router.push(`/products/${targetId}`);
      } catch (err) {
        alert('Failed to save product: ' + (err.response?.data?.message || err.message));
      }
    }
  }
};
</script>

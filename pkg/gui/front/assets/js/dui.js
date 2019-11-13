
(function(l, r) { if (l.getElementById('livereloadscript')) return; r = l.createElement('script'); r.async = 1; r.src = '//' + (window.location.host || 'localhost').split(':')[0] + ':35729/livereload.js?snipver=1'; r.id = 'livereloadscript'; l.head.appendChild(r) })(window.document);
var dui = (function () {
    'use strict';

    function noop() { }
    const identity = x => x;
    function assign(tar, src) {
        // @ts-ignore
        for (const k in src)
            tar[k] = src[k];
        return tar;
    }
    function add_location(element, file, line, column, char) {
        element.__svelte_meta = {
            loc: { file, line, column, char }
        };
    }
    function run(fn) {
        return fn();
    }
    function blank_object() {
        return Object.create(null);
    }
    function run_all(fns) {
        fns.forEach(run);
    }
    function is_function(thing) {
        return typeof thing === 'function';
    }
    function safe_not_equal(a, b) {
        return a != a ? b == b : a !== b || ((a && typeof a === 'object') || typeof a === 'function');
    }
    function validate_store(store, name) {
        if (!store || typeof store.subscribe !== 'function') {
            throw new Error(`'${name}' is not a store with a 'subscribe' method`);
        }
    }
    function subscribe(store, callback) {
        const unsub = store.subscribe(callback);
        return unsub.unsubscribe ? () => unsub.unsubscribe() : unsub;
    }
    function component_subscribe(component, store, callback) {
        component.$$.on_destroy.push(subscribe(store, callback));
    }
    function create_slot(definition, ctx, fn) {
        if (definition) {
            const slot_ctx = get_slot_context(definition, ctx, fn);
            return definition[0](slot_ctx);
        }
    }
    function get_slot_context(definition, ctx, fn) {
        return definition[1]
            ? assign({}, assign(ctx.$$scope.ctx, definition[1](fn ? fn(ctx) : {})))
            : ctx.$$scope.ctx;
    }
    function get_slot_changes(definition, ctx, changed, fn) {
        return definition[1]
            ? assign({}, assign(ctx.$$scope.changed || {}, definition[1](fn ? fn(changed) : {})))
            : ctx.$$scope.changed || {};
    }
    function exclude_internal_props(props) {
        const result = {};
        for (const k in props)
            if (k[0] !== '$')
                result[k] = props[k];
        return result;
    }
    function null_to_empty(value) {
        return value == null ? '' : value;
    }
    const has_prop = (obj, prop) => Object.prototype.hasOwnProperty.call(obj, prop);

    const is_client = typeof window !== 'undefined';
    let now = is_client
        ? () => window.performance.now()
        : () => Date.now();
    let raf = is_client ? cb => requestAnimationFrame(cb) : noop;

    const tasks = new Set();
    let running = false;
    function run_tasks() {
        tasks.forEach(task => {
            if (!task[0](now())) {
                tasks.delete(task);
                task[1]();
            }
        });
        running = tasks.size > 0;
        if (running)
            raf(run_tasks);
    }
    function loop(fn) {
        let task;
        if (!running) {
            running = true;
            raf(run_tasks);
        }
        return {
            promise: new Promise(fulfil => {
                tasks.add(task = [fn, fulfil]);
            }),
            abort() {
                tasks.delete(task);
            }
        };
    }

    function append(target, node) {
        target.appendChild(node);
    }
    function insert(target, node, anchor) {
        target.insertBefore(node, anchor || null);
    }
    function detach(node) {
        node.parentNode.removeChild(node);
    }
    function destroy_each(iterations, detaching) {
        for (let i = 0; i < iterations.length; i += 1) {
            if (iterations[i])
                iterations[i].d(detaching);
        }
    }
    function element(name) {
        return document.createElement(name);
    }
    function svg_element(name) {
        return document.createElementNS('http://www.w3.org/2000/svg', name);
    }
    function text(data) {
        return document.createTextNode(data);
    }
    function space() {
        return text(' ');
    }
    function empty() {
        return text('');
    }
    function listen(node, event, handler, options) {
        node.addEventListener(event, handler, options);
        return () => node.removeEventListener(event, handler, options);
    }
    function attr(node, attribute, value) {
        if (value == null)
            node.removeAttribute(attribute);
        else if (node.getAttribute(attribute) !== value)
            node.setAttribute(attribute, value);
    }
    function set_attributes(node, attributes) {
        // @ts-ignore
        const descriptors = Object.getOwnPropertyDescriptors(node.__proto__);
        for (const key in attributes) {
            if (attributes[key] == null) {
                node.removeAttribute(key);
            }
            else if (key === 'style') {
                node.style.cssText = attributes[key];
            }
            else if (descriptors[key] && descriptors[key].set) {
                node[key] = attributes[key];
            }
            else {
                attr(node, key, attributes[key]);
            }
        }
    }
    function children(element) {
        return Array.from(element.childNodes);
    }
    function set_input_value(input, value) {
        if (value != null || input.value) {
            input.value = value;
        }
    }
    function set_style(node, key, value, important) {
        node.style.setProperty(key, value, important ? 'important' : '');
    }
    function add_resize_listener(element, fn) {
        if (getComputedStyle(element).position === 'static') {
            element.style.position = 'relative';
        }
        const object = document.createElement('object');
        object.setAttribute('style', 'display: block; position: absolute; top: 0; left: 0; height: 100%; width: 100%; overflow: hidden; pointer-events: none; z-index: -1;');
        object.type = 'text/html';
        object.tabIndex = -1;
        let win;
        object.onload = () => {
            win = object.contentDocument.defaultView;
            win.addEventListener('resize', fn);
        };
        if (/Trident/.test(navigator.userAgent)) {
            element.appendChild(object);
            object.data = 'about:blank';
        }
        else {
            object.data = 'about:blank';
            element.appendChild(object);
        }
        return {
            cancel: () => {
                win && win.removeEventListener && win.removeEventListener('resize', fn);
                element.removeChild(object);
            }
        };
    }
    function toggle_class(element, name, toggle) {
        element.classList[toggle ? 'add' : 'remove'](name);
    }
    function custom_event(type, detail) {
        const e = document.createEvent('CustomEvent');
        e.initCustomEvent(type, false, false, detail);
        return e;
    }
    class HtmlTag {
        constructor(html, anchor = null) {
            this.e = element('div');
            this.a = anchor;
            this.u(html);
        }
        m(target, anchor = null) {
            for (let i = 0; i < this.n.length; i += 1) {
                insert(target, this.n[i], anchor);
            }
            this.t = target;
        }
        u(html) {
            this.e.innerHTML = html;
            this.n = Array.from(this.e.childNodes);
        }
        p(html) {
            this.d();
            this.u(html);
            this.m(this.t, this.a);
        }
        d() {
            this.n.forEach(detach);
        }
    }

    let stylesheet;
    let active = 0;
    let current_rules = {};
    // https://github.com/darkskyapp/string-hash/blob/master/index.js
    function hash(str) {
        let hash = 5381;
        let i = str.length;
        while (i--)
            hash = ((hash << 5) - hash) ^ str.charCodeAt(i);
        return hash >>> 0;
    }
    function create_rule(node, a, b, duration, delay, ease, fn, uid = 0) {
        const step = 16.666 / duration;
        let keyframes = '{\n';
        for (let p = 0; p <= 1; p += step) {
            const t = a + (b - a) * ease(p);
            keyframes += p * 100 + `%{${fn(t, 1 - t)}}\n`;
        }
        const rule = keyframes + `100% {${fn(b, 1 - b)}}\n}`;
        const name = `__svelte_${hash(rule)}_${uid}`;
        if (!current_rules[name]) {
            if (!stylesheet) {
                const style = element('style');
                document.head.appendChild(style);
                stylesheet = style.sheet;
            }
            current_rules[name] = true;
            stylesheet.insertRule(`@keyframes ${name} ${rule}`, stylesheet.cssRules.length);
        }
        const animation = node.style.animation || '';
        node.style.animation = `${animation ? `${animation}, ` : ``}${name} ${duration}ms linear ${delay}ms 1 both`;
        active += 1;
        return name;
    }
    function delete_rule(node, name) {
        node.style.animation = (node.style.animation || '')
            .split(', ')
            .filter(name
            ? anim => anim.indexOf(name) < 0 // remove specific animation
            : anim => anim.indexOf('__svelte') === -1 // remove all Svelte animations
        )
            .join(', ');
        if (name && !--active)
            clear_rules();
    }
    function clear_rules() {
        raf(() => {
            if (active)
                return;
            let i = stylesheet.cssRules.length;
            while (i--)
                stylesheet.deleteRule(i);
            current_rules = {};
        });
    }

    let current_component;
    function set_current_component(component) {
        current_component = component;
    }
    function get_current_component() {
        if (!current_component)
            throw new Error(`Function called outside component initialization`);
        return current_component;
    }
    function onMount(fn) {
        get_current_component().$$.on_mount.push(fn);
    }
    function createEventDispatcher() {
        const component = get_current_component();
        return (type, detail) => {
            const callbacks = component.$$.callbacks[type];
            if (callbacks) {
                // TODO are there situations where events could be dispatched
                // in a server (non-DOM) environment?
                const event = custom_event(type, detail);
                callbacks.slice().forEach(fn => {
                    fn.call(component, event);
                });
            }
        };
    }
    // TODO figure out if we still want to support
    // shorthand events, or if we want to implement
    // a real bubbling mechanism
    function bubble(component, event) {
        const callbacks = component.$$.callbacks[event.type];
        if (callbacks) {
            callbacks.slice().forEach(fn => fn(event));
        }
    }

    const dirty_components = [];
    const binding_callbacks = [];
    const render_callbacks = [];
    const flush_callbacks = [];
    const resolved_promise = Promise.resolve();
    let update_scheduled = false;
    function schedule_update() {
        if (!update_scheduled) {
            update_scheduled = true;
            resolved_promise.then(flush);
        }
    }
    function add_render_callback(fn) {
        render_callbacks.push(fn);
    }
    function add_flush_callback(fn) {
        flush_callbacks.push(fn);
    }
    function flush() {
        const seen_callbacks = new Set();
        do {
            // first, call beforeUpdate functions
            // and update components
            while (dirty_components.length) {
                const component = dirty_components.shift();
                set_current_component(component);
                update(component.$$);
            }
            while (binding_callbacks.length)
                binding_callbacks.pop()();
            // then, once components are updated, call
            // afterUpdate functions. This may cause
            // subsequent updates...
            for (let i = 0; i < render_callbacks.length; i += 1) {
                const callback = render_callbacks[i];
                if (!seen_callbacks.has(callback)) {
                    callback();
                    // ...so guard against infinite loops
                    seen_callbacks.add(callback);
                }
            }
            render_callbacks.length = 0;
        } while (dirty_components.length);
        while (flush_callbacks.length) {
            flush_callbacks.pop()();
        }
        update_scheduled = false;
    }
    function update($$) {
        if ($$.fragment !== null) {
            $$.update($$.dirty);
            run_all($$.before_update);
            $$.fragment && $$.fragment.p($$.dirty, $$.ctx);
            $$.dirty = null;
            $$.after_update.forEach(add_render_callback);
        }
    }

    let promise;
    function wait() {
        if (!promise) {
            promise = Promise.resolve();
            promise.then(() => {
                promise = null;
            });
        }
        return promise;
    }
    function dispatch(node, direction, kind) {
        node.dispatchEvent(custom_event(`${direction ? 'intro' : 'outro'}${kind}`));
    }
    const outroing = new Set();
    let outros;
    function group_outros() {
        outros = {
            r: 0,
            c: [],
            p: outros // parent group
        };
    }
    function check_outros() {
        if (!outros.r) {
            run_all(outros.c);
        }
        outros = outros.p;
    }
    function transition_in(block, local) {
        if (block && block.i) {
            outroing.delete(block);
            block.i(local);
        }
    }
    function transition_out(block, local, detach, callback) {
        if (block && block.o) {
            if (outroing.has(block))
                return;
            outroing.add(block);
            outros.c.push(() => {
                outroing.delete(block);
                if (callback) {
                    if (detach)
                        block.d(1);
                    callback();
                }
            });
            block.o(local);
        }
    }
    const null_transition = { duration: 0 };
    function create_in_transition(node, fn, params) {
        let config = fn(node, params);
        let running = false;
        let animation_name;
        let task;
        let uid = 0;
        function cleanup() {
            if (animation_name)
                delete_rule(node, animation_name);
        }
        function go() {
            const { delay = 0, duration = 300, easing = identity, tick = noop, css } = config || null_transition;
            if (css)
                animation_name = create_rule(node, 0, 1, duration, delay, easing, css, uid++);
            tick(0, 1);
            const start_time = now() + delay;
            const end_time = start_time + duration;
            if (task)
                task.abort();
            running = true;
            add_render_callback(() => dispatch(node, true, 'start'));
            task = loop(now => {
                if (running) {
                    if (now >= end_time) {
                        tick(1, 0);
                        dispatch(node, true, 'end');
                        cleanup();
                        return running = false;
                    }
                    if (now >= start_time) {
                        const t = easing((now - start_time) / duration);
                        tick(t, 1 - t);
                    }
                }
                return running;
            });
        }
        let started = false;
        return {
            start() {
                if (started)
                    return;
                delete_rule(node);
                if (is_function(config)) {
                    config = config();
                    wait().then(go);
                }
                else {
                    go();
                }
            },
            invalidate() {
                started = false;
            },
            end() {
                if (running) {
                    cleanup();
                    running = false;
                }
            }
        };
    }
    function create_out_transition(node, fn, params) {
        let config = fn(node, params);
        let running = true;
        let animation_name;
        const group = outros;
        group.r += 1;
        function go() {
            const { delay = 0, duration = 300, easing = identity, tick = noop, css } = config || null_transition;
            if (css)
                animation_name = create_rule(node, 1, 0, duration, delay, easing, css);
            const start_time = now() + delay;
            const end_time = start_time + duration;
            add_render_callback(() => dispatch(node, false, 'start'));
            loop(now => {
                if (running) {
                    if (now >= end_time) {
                        tick(0, 1);
                        dispatch(node, false, 'end');
                        if (!--group.r) {
                            // this will result in `end()` being called,
                            // so we don't need to clean up here
                            run_all(group.c);
                        }
                        return false;
                    }
                    if (now >= start_time) {
                        const t = easing((now - start_time) / duration);
                        tick(1 - t, t);
                    }
                }
                return running;
            });
        }
        if (is_function(config)) {
            wait().then(() => {
                // @ts-ignore
                config = config();
                go();
            });
        }
        else {
            go();
        }
        return {
            end(reset) {
                if (reset && config.tick) {
                    config.tick(1, 0);
                }
                if (running) {
                    if (animation_name)
                        delete_rule(node, animation_name);
                    running = false;
                }
            }
        };
    }
    function create_bidirectional_transition(node, fn, params, intro) {
        let config = fn(node, params);
        let t = intro ? 0 : 1;
        let running_program = null;
        let pending_program = null;
        let animation_name = null;
        function clear_animation() {
            if (animation_name)
                delete_rule(node, animation_name);
        }
        function init(program, duration) {
            const d = program.b - t;
            duration *= Math.abs(d);
            return {
                a: t,
                b: program.b,
                d,
                duration,
                start: program.start,
                end: program.start + duration,
                group: program.group
            };
        }
        function go(b) {
            const { delay = 0, duration = 300, easing = identity, tick = noop, css } = config || null_transition;
            const program = {
                start: now() + delay,
                b
            };
            if (!b) {
                // @ts-ignore todo: improve typings
                program.group = outros;
                outros.r += 1;
            }
            if (running_program) {
                pending_program = program;
            }
            else {
                // if this is an intro, and there's a delay, we need to do
                // an initial tick and/or apply CSS animation immediately
                if (css) {
                    clear_animation();
                    animation_name = create_rule(node, t, b, duration, delay, easing, css);
                }
                if (b)
                    tick(0, 1);
                running_program = init(program, duration);
                add_render_callback(() => dispatch(node, b, 'start'));
                loop(now => {
                    if (pending_program && now > pending_program.start) {
                        running_program = init(pending_program, duration);
                        pending_program = null;
                        dispatch(node, running_program.b, 'start');
                        if (css) {
                            clear_animation();
                            animation_name = create_rule(node, t, running_program.b, running_program.duration, 0, easing, config.css);
                        }
                    }
                    if (running_program) {
                        if (now >= running_program.end) {
                            tick(t = running_program.b, 1 - t);
                            dispatch(node, running_program.b, 'end');
                            if (!pending_program) {
                                // we're done
                                if (running_program.b) {
                                    // intro — we can tidy up immediately
                                    clear_animation();
                                }
                                else {
                                    // outro — needs to be coordinated
                                    if (!--running_program.group.r)
                                        run_all(running_program.group.c);
                                }
                            }
                            running_program = null;
                        }
                        else if (now >= running_program.start) {
                            const p = now - running_program.start;
                            t = running_program.a + running_program.d * easing(p / running_program.duration);
                            tick(t, 1 - t);
                        }
                    }
                    return !!(running_program || pending_program);
                });
            }
        }
        return {
            run(b) {
                if (is_function(config)) {
                    wait().then(() => {
                        // @ts-ignore
                        config = config();
                        go(b);
                    });
                }
                else {
                    go(b);
                }
            },
            end() {
                clear_animation();
                running_program = pending_program = null;
            }
        };
    }

    const globals = (typeof window !== 'undefined' ? window : global);

    function get_spread_update(levels, updates) {
        const update = {};
        const to_null_out = {};
        const accounted_for = { $$scope: 1 };
        let i = levels.length;
        while (i--) {
            const o = levels[i];
            const n = updates[i];
            if (n) {
                for (const key in o) {
                    if (!(key in n))
                        to_null_out[key] = 1;
                }
                for (const key in n) {
                    if (!accounted_for[key]) {
                        update[key] = n[key];
                        accounted_for[key] = 1;
                    }
                }
                levels[i] = n;
            }
            else {
                for (const key in o) {
                    accounted_for[key] = 1;
                }
            }
        }
        for (const key in to_null_out) {
            if (!(key in update))
                update[key] = undefined;
        }
        return update;
    }
    function get_spread_object(spread_props) {
        return typeof spread_props === 'object' && spread_props !== null ? spread_props : {};
    }

    function bind(component, name, callback) {
        if (has_prop(component.$$.props, name)) {
            name = component.$$.props[name] || name;
            component.$$.bound[name] = callback;
            callback(component.$$.ctx[name]);
        }
    }
    function create_component(block) {
        block && block.c();
    }
    function mount_component(component, target, anchor) {
        const { fragment, on_mount, on_destroy, after_update } = component.$$;
        fragment && fragment.m(target, anchor);
        // onMount happens before the initial afterUpdate
        add_render_callback(() => {
            const new_on_destroy = on_mount.map(run).filter(is_function);
            if (on_destroy) {
                on_destroy.push(...new_on_destroy);
            }
            else {
                // Edge case - component was destroyed immediately,
                // most likely as a result of a binding initialising
                run_all(new_on_destroy);
            }
            component.$$.on_mount = [];
        });
        after_update.forEach(add_render_callback);
    }
    function destroy_component(component, detaching) {
        const $$ = component.$$;
        if ($$.fragment !== null) {
            run_all($$.on_destroy);
            $$.fragment && $$.fragment.d(detaching);
            // TODO null out other refs, including component.$$ (but need to
            // preserve final state?)
            $$.on_destroy = $$.fragment = null;
            $$.ctx = {};
        }
    }
    function make_dirty(component, key) {
        if (!component.$$.dirty) {
            dirty_components.push(component);
            schedule_update();
            component.$$.dirty = blank_object();
        }
        component.$$.dirty[key] = true;
    }
    function init(component, options, instance, create_fragment, not_equal, props) {
        const parent_component = current_component;
        set_current_component(component);
        const prop_values = options.props || {};
        const $$ = component.$$ = {
            fragment: null,
            ctx: null,
            // state
            props,
            update: noop,
            not_equal,
            bound: blank_object(),
            // lifecycle
            on_mount: [],
            on_destroy: [],
            before_update: [],
            after_update: [],
            context: new Map(parent_component ? parent_component.$$.context : []),
            // everything else
            callbacks: blank_object(),
            dirty: null
        };
        let ready = false;
        $$.ctx = instance
            ? instance(component, prop_values, (key, ret, value = ret) => {
                if ($$.ctx && not_equal($$.ctx[key], $$.ctx[key] = value)) {
                    if ($$.bound[key])
                        $$.bound[key](value);
                    if (ready)
                        make_dirty(component, key);
                }
                return ret;
            })
            : prop_values;
        $$.update();
        ready = true;
        run_all($$.before_update);
        // `false` as a special case of no DOM component
        $$.fragment = create_fragment ? create_fragment($$.ctx) : false;
        if (options.target) {
            if (options.hydrate) {
                // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
                $$.fragment && $$.fragment.l(children(options.target));
            }
            else {
                // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
                $$.fragment && $$.fragment.c();
            }
            if (options.intro)
                transition_in(component.$$.fragment);
            mount_component(component, options.target, options.anchor);
            flush();
        }
        set_current_component(parent_component);
    }
    class SvelteComponent {
        $destroy() {
            destroy_component(this, 1);
            this.$destroy = noop;
        }
        $on(type, callback) {
            const callbacks = (this.$$.callbacks[type] || (this.$$.callbacks[type] = []));
            callbacks.push(callback);
            return () => {
                const index = callbacks.indexOf(callback);
                if (index !== -1)
                    callbacks.splice(index, 1);
            };
        }
        $set() {
            // overridden by instance, if it has props
        }
    }

    function dispatch_dev(type, detail) {
        document.dispatchEvent(custom_event(type, detail));
    }
    function append_dev(target, node) {
        dispatch_dev("SvelteDOMInsert", { target, node });
        append(target, node);
    }
    function insert_dev(target, node, anchor) {
        dispatch_dev("SvelteDOMInsert", { target, node, anchor });
        insert(target, node, anchor);
    }
    function detach_dev(node) {
        dispatch_dev("SvelteDOMRemove", { node });
        detach(node);
    }
    function listen_dev(node, event, handler, options, has_prevent_default, has_stop_propagation) {
        const modifiers = options === true ? ["capture"] : options ? Array.from(Object.keys(options)) : [];
        if (has_prevent_default)
            modifiers.push('preventDefault');
        if (has_stop_propagation)
            modifiers.push('stopPropagation');
        dispatch_dev("SvelteDOMAddEventListener", { node, event, handler, modifiers });
        const dispose = listen(node, event, handler, options);
        return () => {
            dispatch_dev("SvelteDOMRemoveEventListener", { node, event, handler, modifiers });
            dispose();
        };
    }
    function attr_dev(node, attribute, value) {
        attr(node, attribute, value);
        if (value == null)
            dispatch_dev("SvelteDOMRemoveAttribute", { node, attribute });
        else
            dispatch_dev("SvelteDOMSetAttribute", { node, attribute, value });
    }
    function prop_dev(node, property, value) {
        node[property] = value;
        dispatch_dev("SvelteDOMSetProperty", { node, property, value });
    }
    function set_data_dev(text, data) {
        data = '' + data;
        if (text.data === data)
            return;
        dispatch_dev("SvelteDOMSetData", { node: text, data });
        text.data = data;
    }
    class SvelteComponentDev extends SvelteComponent {
        constructor(options) {
            if (!options || (!options.target && !options.$$inline)) {
                throw new Error(`'target' is a required option`);
            }
            super();
        }
        $destroy() {
            super.$destroy();
            this.$destroy = () => {
                console.warn(`Component was already destroyed`); // eslint-disable-line no-console
            };
        }
    }

    function cubicInOut(t) {
        return t < 0.5 ? 4.0 * t * t * t : 0.5 * Math.pow(2.0 * t - 2.0, 3.0) + 1.0;
    }
    function cubicOut(t) {
        const f = t - 1.0;
        return f * f * f + 1.0;
    }
    function quadOut(t) {
        return -t * (t - 2.0);
    }

    function fade(node, { delay = 0, duration = 400, easing = identity }) {
        const o = +getComputedStyle(node).opacity;
        return {
            delay,
            duration,
            easing,
            css: t => `opacity: ${t * o}`
        };
    }
    function fly(node, { delay = 0, duration = 400, easing = cubicOut, x = 0, y = 0, opacity = 0 }) {
        const style = getComputedStyle(node);
        const target_opacity = +style.opacity;
        const transform = style.transform === 'none' ? '' : style.transform;
        const od = target_opacity * (1 - opacity);
        return {
            delay,
            duration,
            easing,
            css: (t, u) => `
			transform: ${transform} translate(${(1 - t) * x}px, ${(1 - t) * y}px);
			opacity: ${target_opacity - (od * u)}`
        };
    }
    function slide(node, { delay = 0, duration = 400, easing = cubicOut }) {
        const style = getComputedStyle(node);
        const opacity = +style.opacity;
        const height = parseFloat(style.height);
        const padding_top = parseFloat(style.paddingTop);
        const padding_bottom = parseFloat(style.paddingBottom);
        const margin_top = parseFloat(style.marginTop);
        const margin_bottom = parseFloat(style.marginBottom);
        const border_top_width = parseFloat(style.borderTopWidth);
        const border_bottom_width = parseFloat(style.borderBottomWidth);
        return {
            delay,
            duration,
            easing,
            css: t => `overflow: hidden;` +
                `opacity: ${Math.min(t * 20, 1) * opacity};` +
                `height: ${t * height}px;` +
                `padding-top: ${t * padding_top}px;` +
                `padding-bottom: ${t * padding_bottom}px;` +
                `margin-top: ${t * margin_top}px;` +
                `margin-bottom: ${t * margin_bottom}px;` +
                `border-top-width: ${t * border_top_width}px;` +
                `border-bottom-width: ${t * border_bottom_width}px;`
        };
    }
    function draw(node, { delay = 0, speed, duration, easing = cubicInOut }) {
        const len = node.getTotalLength();
        if (duration === undefined) {
            if (speed === undefined) {
                duration = 800;
            }
            else {
                duration = len / speed;
            }
        }
        else if (typeof duration === 'function') {
            duration = duration(len);
        }
        return {
            delay,
            duration,
            easing,
            css: (t, u) => `stroke-dasharray: ${t * len} ${u * len}`
        };
    }

    const inner = `M77.08,2.55c3.87,1.03 6.96,2.58 10.32,4.64c5.93,3.87 10.58,8.51 14.19,14.71c3.87,6.19 5.42,13.16 5.42,20.64c0,7.22 -1.81,14.18 -5.41,20.37c-3.61,6.45 -8.25,11.35 -14.19,14.96c-3.35,2.06 -6.96,3.87 -10.32,4.9c-3.87,1.03 -7.74,1.55 -11.61,1.55v-14.45c6.96,-0.26 13.42,-2.58 19.09,-8c5.67,-5.42 8.51,-11.87 8.51,-19.61c0,-7.74 -2.58,-14.19 -7.99,-19.6c-5.42,-5.42 -11.86,-8 -19.6,-8c-7.74,0 -14.44,2.58 -19.6,8c-5.42,5.42 -8,11.87 -8,19.6l0,85.9c-3.1,-3.1 -7.99,-7.74 -13.93,-13.67v-72.23c0,-3.87 0.52,-7.73 1.55,-11.35c1.03,-3.87 2.58,-7.22 4.64,-10.32c3.87,-5.93 8.52,-10.58 14.71,-14.45c6.19,-3.61 13.16,-5.16 20.64,-5.16c3.87,0 8,0.52 11.61,1.55zM78.37,42.28c0,7.22 -5.93,13.16 -13.15,13.16c-7.48,0.26 -13.16,-5.68 -13.16,-13.16c0,-7.22 5.94,-13.16 13.16,-13.16c7.22,0 13.15,5.93 13.15,13.16zM13.63,37.12l0,69.39c-6.19,-6.19 -11.09,-10.83 -13.93,-13.93l0,-55.46z`;

    /* src/boot/logo/BootLogo.svelte generated by Svelte v3.14.0 */
    const file = "src/boot/logo/BootLogo.svelte";

    function get_each_context(ctx, list, i) {
    	const child_ctx = Object.create(ctx);
    	child_ctx.char = list[i];
    	child_ctx.i = i;
    	return child_ctx;
    }

    function get_each_context_1(ctx, list, i) {
    	const child_ctx = Object.create(ctx);
    	child_ctx.char = list[i];
    	child_ctx.i = i;
    	return child_ctx;
    }

    // (27:2) {#each 'PLAN 9 FROM CRYPTO SPACE' as char, i}
    function create_each_block_1(ctx) {
    	let span;
    	let t;
    	let span_intro;

    	const block = {
    		c: function create() {
    			span = element("span");
    			t = text(ctx.char);
    			attr_dev(span, "class", "svelte-rpr53z");
    			add_location(span, file, 27, 3, 570);
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, span, anchor);
    			append_dev(span, t);
    		},
    		p: noop,
    		i: function intro(local) {
    			if (!span_intro) {
    				add_render_callback(() => {
    					span_intro = create_in_transition(span, fade, { delay: 33 + ctx.i * 150, duration: 3333 });
    					span_intro.start();
    				});
    			}
    		},
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(span);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_each_block_1.name,
    		type: "each",
    		source: "(27:2) {#each 'PLAN 9 FROM CRYPTO SPACE' as char, i}",
    		ctx
    	});

    	return block;
    }

    // (44:2) {#each 'ParallelCoin' as char, i}
    function create_each_block(ctx) {
    	let span;
    	let t;
    	let span_intro;

    	const block = {
    		c: function create() {
    			span = element("span");
    			t = text(ctx.char);
    			attr_dev(span, "class", "svelte-rpr53z");
    			add_location(span, file, 44, 3, 1008);
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, span, anchor);
    			append_dev(span, t);
    		},
    		p: noop,
    		i: function intro(local) {
    			if (!span_intro) {
    				add_render_callback(() => {
    					span_intro = create_in_transition(span, fade, { delay: 3333 + ctx.i * 150, duration: 999 });
    					span_intro.start();
    				});
    			}
    		},
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(span);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_each_block.name,
    		type: "each",
    		source: "(44:2) {#each 'ParallelCoin' as char, i}",
    		ctx
    	});

    	return block;
    }

    function create_fragment(ctx) {
    	let div0;
    	let div0_outro;
    	let t0;
    	let svg;
    	let g;
    	let path;
    	let path_intro;
    	let g_outro;
    	let t1;
    	let div1;
    	let div1_outro;
    	let t2;
    	let div2;
    	let caption;
    	let t3;
    	let t4;
    	let current;
    	let each_value_1 = "PLAN 9 FROM CRYPTO SPACE";
    	let each_blocks_1 = [];

    	for (let i = 0; i < each_value_1.length; i += 1) {
    		each_blocks_1[i] = create_each_block_1(get_each_context_1(ctx, each_value_1, i));
    	}

    	let each_value = "ParallelCoin";
    	let each_blocks = [];

    	for (let i = 0; i < each_value.length; i += 1) {
    		each_blocks[i] = create_each_block(get_each_context(ctx, each_value, i));
    	}

    	const block = {
    		c: function create() {
    			div0 = element("div");

    			for (let i = 0; i < each_blocks_1.length; i += 1) {
    				each_blocks_1[i].c();
    			}

    			t0 = space();
    			svg = svg_element("svg");
    			g = svg_element("g");
    			path = svg_element("path");
    			t1 = space();
    			div1 = element("div");

    			for (let i = 0; i < each_blocks.length; i += 1) {
    				each_blocks[i].c();
    			}

    			t2 = space();
    			div2 = element("div");
    			caption = element("caption");
    			t3 = text(ctx.progress);
    			t4 = text("%");
    			attr_dev(div0, "class", "centered marginTop plan svelte-rpr53z");
    			add_location(div0, file, 25, 0, 445);
    			set_style(path, "stroke", "#cfcfcf");
    			set_style(path, "stroke-width", "1.5");
    			attr_dev(path, "d", inner);
    			attr_dev(path, "class", "svelte-rpr53z");
    			add_location(path, file, 35, 3, 783);
    			attr_dev(g, "opacity", "0.2");
    			add_location(g, file, 34, 2, 735);
    			attr_dev(svg, "id", "bootlogo");
    			attr_dev(svg, "class", "marginTopBig svelte-rpr53z");
    			attr_dev(svg, "viewBox", "0 0 108 128");
    			add_location(svg, file, 33, 1, 669);
    			attr_dev(div1, "class", "centered name svelte-rpr53z");
    			add_location(div1, file, 42, 1, 904);
    			add_location(caption, file, 53, 1, 1164);
    			attr_dev(div2, "class", "progress justifyCenter textCenter txDark svelte-rpr53z");
    			add_location(div2, file, 52, 0, 1108);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div0, anchor);

    			for (let i = 0; i < each_blocks_1.length; i += 1) {
    				each_blocks_1[i].m(div0, null);
    			}

    			insert_dev(target, t0, anchor);
    			insert_dev(target, svg, anchor);
    			append_dev(svg, g);
    			append_dev(g, path);
    			insert_dev(target, t1, anchor);
    			insert_dev(target, div1, anchor);

    			for (let i = 0; i < each_blocks.length; i += 1) {
    				each_blocks[i].m(div1, null);
    			}

    			insert_dev(target, t2, anchor);
    			insert_dev(target, div2, anchor);
    			append_dev(div2, caption);
    			append_dev(caption, t3);
    			append_dev(caption, t4);
    			current = true;
    		},
    		p: function update(changed, ctx) {
    			if (!current || changed.progress) set_data_dev(t3, ctx.progress);
    		},
    		i: function intro(local) {
    			if (current) return;

    			for (let i = 0; i < each_value_1.length; i += 1) {
    				transition_in(each_blocks_1[i]);
    			}

    			if (div0_outro) div0_outro.end(1);

    			if (!path_intro) {
    				add_render_callback(() => {
    					path_intro = create_in_transition(path, draw, { duration: 3333 });
    					path_intro.start();
    				});
    			}

    			if (g_outro) g_outro.end(1);

    			for (let i = 0; i < each_value.length; i += 1) {
    				transition_in(each_blocks[i]);
    			}

    			if (div1_outro) div1_outro.end(1);
    			current = true;
    		},
    		o: function outro(local) {
    			div0_outro = create_out_transition(div0, fly, { y: -40, duration: 999 });
    			g_outro = create_out_transition(g, fade, { duration: 999 });
    			div1_outro = create_out_transition(div1, fly, { y: -20, duration: 9999 });
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div0);
    			destroy_each(each_blocks_1, detaching);
    			if (detaching && div0_outro) div0_outro.end();
    			if (detaching) detach_dev(t0);
    			if (detaching) detach_dev(svg);
    			if (detaching && g_outro) g_outro.end();
    			if (detaching) detach_dev(t1);
    			if (detaching) detach_dev(div1);
    			destroy_each(each_blocks, detaching);
    			if (detaching && div1_outro) div1_outro.end();
    			if (detaching) detach_dev(t2);
    			if (detaching) detach_dev(div2);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    function instance($$self, $$props, $$invalidate) {
    	let progress = 0;

    	function next() {
    		setTimeout(
    			() => {
    				if (progress === 100) {
    					$$invalidate("progress", progress = 0);
    				}

    				$$invalidate("progress", progress += 1);
    				next();
    			},
    			100
    		);
    	}

    	next();

    	$$self.$capture_state = () => {
    		return {};
    	};

    	$$self.$inject_state = $$props => {
    		if ("progress" in $$props) $$invalidate("progress", progress = $$props.progress);
    	};

    	return { progress };
    }

    class BootLogo extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, instance, create_fragment, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "BootLogo",
    			options,
    			id: create_fragment.name
    		});
    	}
    }

    const subscriber_queue = [];
    /**
     * Creates a `Readable` store that allows reading by subscription.
     * @param value initial value
     * @param {StartStopNotifier}start start and stop notifications for subscriptions
     */
    function readable(value, start) {
        return {
            subscribe: writable(value, start).subscribe,
        };
    }
    /**
     * Create a `Writable` store that allows both updating and reading by subscription.
     * @param {*=}value initial value
     * @param {StartStopNotifier=}start start and stop notifications for subscriptions
     */
    function writable(value, start = noop) {
        let stop;
        const subscribers = [];
        function set(new_value) {
            if (safe_not_equal(value, new_value)) {
                value = new_value;
                if (stop) { // store is ready
                    const run_queue = !subscriber_queue.length;
                    for (let i = 0; i < subscribers.length; i += 1) {
                        const s = subscribers[i];
                        s[1]();
                        subscriber_queue.push(s, value);
                    }
                    if (run_queue) {
                        for (let i = 0; i < subscriber_queue.length; i += 2) {
                            subscriber_queue[i][0](subscriber_queue[i + 1]);
                        }
                        subscriber_queue.length = 0;
                    }
                }
            }
        }
        function update(fn) {
            set(fn(value));
        }
        function subscribe(run, invalidate = noop) {
            const subscriber = [run, invalidate];
            subscribers.push(subscriber);
            if (subscribers.length === 1) {
                stop = start(set) || noop;
            }
            run(value);
            return () => {
                const index = subscribers.indexOf(subscriber);
                if (index !== -1) {
                    subscribers.splice(index, 1);
                }
                if (subscribers.length === 0) {
                    stop();
                    stop = null;
                }
            };
        }
        return { set, update, subscribe };
    }

    const bios = readable([], function start(set) {
        const interval = setInterval(() => {
            fetch(`http://127.0.0.1:3999/bios`)
                .then(resp => resp.json())
                .then(data => {
                    set(data);
                });
        }, 100);
        return function stop() {
            clearInterval(interval);
        };
    });

    /* src/boot/Boot.svelte generated by Svelte v3.14.0 */
    const file$1 = "src/boot/Boot.svelte";

    // (9:0) {#if bios.logo}
    function create_if_block(ctx) {
    	let current;
    	const bootlogo = new BootLogo({ $$inline: true });

    	const block = {
    		c: function create() {
    			create_component(bootlogo.$$.fragment);
    		},
    		m: function mount(target, anchor) {
    			mount_component(bootlogo, target, anchor);
    			current = true;
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(bootlogo.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(bootlogo.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			destroy_component(bootlogo, detaching);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_if_block.name,
    		type: "if",
    		source: "(9:0) {#if bios.logo}",
    		ctx
    	});

    	return block;
    }

    function create_fragment$1(ctx) {
    	let div;
    	let h1;
    	let t1;
    	let current;
    	let if_block = bios.logo && create_if_block(ctx);

    	const block = {
    		c: function create() {
    			div = element("div");
    			h1 = element("h1");
    			h1.textContent = "Boot";
    			t1 = space();
    			if (if_block) if_block.c();
    			attr_dev(h1, "id", "biosMessages");
    			add_location(h1, file$1, 6, 0, 192);
    			attr_dev(div, "class", "fullScreen flx flc justifyEvenly itemsCenter bgDark boot");
    			add_location(div, file$1, 5, 0, 121);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div, anchor);
    			append_dev(div, h1);
    			append_dev(div, t1);
    			if (if_block) if_block.m(div, null);
    			current = true;
    		},
    		p: noop,
    		i: function intro(local) {
    			if (current) return;
    			transition_in(if_block);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(if_block);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div);
    			if (if_block) if_block.d();
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$1.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class Boot extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$1, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "Boot",
    			options,
    			id: create_fragment$1.name
    		});
    	}
    }

    /* src/components/panels/PanelBalance.svelte generated by Svelte v3.14.0 */

    const file$2 = "src/components/panels/PanelBalance.svelte";

    function create_fragment$2(ctx) {
    	let div7;
    	let div6;
    	let div4;
    	let div2;
    	let div0;
    	let t1;
    	let div1;
    	let span0;
    	let t2;
    	let t3;
    	let div3;
    	let t4;
    	let div5;
    	let small0;
    	let span1;
    	let strong0;
    	let span2;
    	let t6;
    	let small1;
    	let span3;
    	let strong1;
    	let span4;

    	const block = {
    		c: function create() {
    			div7 = element("div");
    			div6 = element("div");
    			div4 = element("div");
    			div2 = element("div");
    			div0 = element("div");
    			div0.textContent = "Balance:";
    			t1 = space();
    			div1 = element("div");
    			span0 = element("span");
    			t2 = text(" DUO");
    			t3 = space();
    			div3 = element("div");
    			t4 = space();
    			div5 = element("div");
    			small0 = element("small");
    			span1 = element("span");
    			span1.textContent = "Pending: ";
    			strong0 = element("strong");
    			span2 = element("span");
    			t6 = space();
    			small1 = element("small");
    			span3 = element("span");
    			span3.textContent = "Transactions: ";
    			strong1 = element("strong");
    			span4 = element("span");
    			attr_dev(div0, "class", "e-card-header-title");
    			add_location(div0, file$2, 33, 4, 1727);
    			attr_dev(span0, "v-html", "this.duOSys.status.balance.balance");
    			attr_dev(span0, "class", "svelte-f3witw");
    			add_location(span0, file$2, 34, 34, 1809);
    			attr_dev(div1, "class", "e-card-sub-title svelte-f3witw");
    			add_location(div1, file$2, 34, 4, 1779);
    			attr_dev(div2, "class", "e-card-header-caption");
    			add_location(div2, file$2, 32, 3, 1687);
    			attr_dev(div3, "class", "e-card-header-image balance svelte-f3witw");
    			add_location(div3, file$2, 36, 3, 1890);
    			attr_dev(div4, "class", "e-card-header");
    			add_location(div4, file$2, 31, 2, 1656);
    			attr_dev(span1, "class", "svelte-f3witw");
    			add_location(span1, file$2, 39, 10, 1996);
    			attr_dev(span2, "v-html", "this.duOSys.status.balance.unconfirmed");
    			attr_dev(span2, "class", "svelte-f3witw");
    			add_location(span2, file$2, 39, 40, 2026);
    			add_location(strong0, file$2, 39, 32, 2018);
    			add_location(small0, file$2, 39, 3, 1989);
    			attr_dev(span3, "class", "svelte-f3witw");
    			add_location(span3, file$2, 40, 10, 2115);
    			attr_dev(span4, "v-html", "this.duOSys.status.txsnumber");
    			attr_dev(span4, "class", "svelte-f3witw");
    			add_location(span4, file$2, 40, 45, 2150);
    			add_location(strong1, file$2, 40, 37, 2142);
    			add_location(small1, file$2, 40, 3, 2108);
    			attr_dev(div5, "class", "flx flc e-card-content svelte-f3witw");
    			add_location(div5, file$2, 38, 2, 1949);
    			attr_dev(div6, "class", "e-card flx flc justifyBetween duoCard svelte-f3witw");
    			add_location(div6, file$2, 30, 1, 1602);
    			attr_dev(div7, "class", "rwrap flx");
    			add_location(div7, file$2, 29, 0, 1577);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div7, anchor);
    			append_dev(div7, div6);
    			append_dev(div6, div4);
    			append_dev(div4, div2);
    			append_dev(div2, div0);
    			append_dev(div2, t1);
    			append_dev(div2, div1);
    			append_dev(div1, span0);
    			append_dev(div1, t2);
    			append_dev(div4, t3);
    			append_dev(div4, div3);
    			append_dev(div6, t4);
    			append_dev(div6, div5);
    			append_dev(div5, small0);
    			append_dev(small0, span1);
    			append_dev(small0, strong0);
    			append_dev(strong0, span2);
    			append_dev(div5, t6);
    			append_dev(div5, small1);
    			append_dev(small1, span3);
    			append_dev(small1, strong1);
    			append_dev(strong1, span4);
    		},
    		p: noop,
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div7);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$2.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class PanelBalance extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$2, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PanelBalance",
    			options,
    			id: create_fragment$2.name
    		});
    	}
    }

    /* node_modules/smelte/src/components/Icon/Icon.svelte generated by Svelte v3.14.0 */

    const file$3 = "node_modules/smelte/src/components/Icon/Icon.svelte";

    function create_fragment$3(ctx) {
    	let i;
    	let i_class_value;
    	let current;
    	const default_slot_template = ctx.$$slots.default;
    	const default_slot = create_slot(default_slot_template, ctx, null);

    	const block = {
    		c: function create() {
    			i = element("i");
    			if (default_slot) default_slot.c();
    			attr_dev(i, "aria-hidden", "true");
    			attr_dev(i, "class", i_class_value = "material-icons " + ctx.className + " transition" + " svelte-1d8l9tj");
    			set_style(i, "color", ctx.color);
    			toggle_class(i, "reverse", ctx.reverse);
    			toggle_class(i, "tip", ctx.tip);
    			toggle_class(i, "text-base", ctx.small);
    			toggle_class(i, "text-xs", ctx.xs);
    			add_location(i, file$3, 24, 0, 1012);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, i, anchor);

    			if (default_slot) {
    				default_slot.m(i, null);
    			}

    			current = true;
    		},
    		p: function update(changed, ctx) {
    			if (default_slot && default_slot.p && changed.$$scope) {
    				default_slot.p(get_slot_changes(default_slot_template, ctx, changed, null), get_slot_context(default_slot_template, ctx, null));
    			}

    			if (!current || changed.className && i_class_value !== (i_class_value = "material-icons " + ctx.className + " transition" + " svelte-1d8l9tj")) {
    				attr_dev(i, "class", i_class_value);
    			}

    			if (!current || changed.color) {
    				set_style(i, "color", ctx.color);
    			}

    			if (changed.className || changed.reverse) {
    				toggle_class(i, "reverse", ctx.reverse);
    			}

    			if (changed.className || changed.tip) {
    				toggle_class(i, "tip", ctx.tip);
    			}

    			if (changed.className || changed.small) {
    				toggle_class(i, "text-base", ctx.small);
    			}

    			if (changed.className || changed.xs) {
    				toggle_class(i, "text-xs", ctx.xs);
    			}
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(default_slot, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(default_slot, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(i);
    			if (default_slot) default_slot.d(detaching);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$3.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    function instance$1($$self, $$props, $$invalidate) {
    	let { class: className = "" } = $$props;
    	let { small = false } = $$props;
    	let { xs = false } = $$props;
    	let { reverse = false } = $$props;
    	let { tip = false } = $$props;
    	let { color = "" } = $$props;
    	const writable_props = ["class", "small", "xs", "reverse", "tip", "color"];

    	Object.keys($$props).forEach(key => {
    		if (!writable_props.includes(key) && !key.startsWith("$$")) console.warn(`<Icon> was created with unknown prop '${key}'`);
    	});

    	let { $$slots = {}, $$scope } = $$props;

    	$$self.$set = $$props => {
    		if ("class" in $$props) $$invalidate("className", className = $$props.class);
    		if ("small" in $$props) $$invalidate("small", small = $$props.small);
    		if ("xs" in $$props) $$invalidate("xs", xs = $$props.xs);
    		if ("reverse" in $$props) $$invalidate("reverse", reverse = $$props.reverse);
    		if ("tip" in $$props) $$invalidate("tip", tip = $$props.tip);
    		if ("color" in $$props) $$invalidate("color", color = $$props.color);
    		if ("$$scope" in $$props) $$invalidate("$$scope", $$scope = $$props.$$scope);
    	};

    	$$self.$capture_state = () => {
    		return {
    			className,
    			small,
    			xs,
    			reverse,
    			tip,
    			color
    		};
    	};

    	$$self.$inject_state = $$props => {
    		if ("className" in $$props) $$invalidate("className", className = $$props.className);
    		if ("small" in $$props) $$invalidate("small", small = $$props.small);
    		if ("xs" in $$props) $$invalidate("xs", xs = $$props.xs);
    		if ("reverse" in $$props) $$invalidate("reverse", reverse = $$props.reverse);
    		if ("tip" in $$props) $$invalidate("tip", tip = $$props.tip);
    		if ("color" in $$props) $$invalidate("color", color = $$props.color);
    	};

    	return {
    		className,
    		small,
    		xs,
    		reverse,
    		tip,
    		color,
    		$$slots,
    		$$scope
    	};
    }

    class Icon extends SvelteComponentDev {
    	constructor(options) {
    		super(options);

    		init(this, options, instance$1, create_fragment$3, safe_not_equal, {
    			class: "className",
    			small: 0,
    			xs: 0,
    			reverse: 0,
    			tip: 0,
    			color: 0
    		});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "Icon",
    			options,
    			id: create_fragment$3.name
    		});
    	}

    	get class() {
    		throw new Error("<Icon>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set class(value) {
    		throw new Error("<Icon>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get small() {
    		throw new Error("<Icon>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set small(value) {
    		throw new Error("<Icon>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get xs() {
    		throw new Error("<Icon>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set xs(value) {
    		throw new Error("<Icon>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get reverse() {
    		throw new Error("<Icon>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set reverse(value) {
    		throw new Error("<Icon>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get tip() {
    		throw new Error("<Icon>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set tip(value) {
    		throw new Error("<Icon>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get color() {
    		throw new Error("<Icon>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set color(value) {
    		throw new Error("<Icon>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}
    }

    const noDepth = ["white", "black", "transparent"];

    function getClass(prop, color, depth, defaultDepth) {
      if (noDepth.includes(color)) {
        return `${prop}-${color}`;
      }
      return `${prop}-${color}-${depth || defaultDepth} `;
    }

    function utils(color, defaultDepth = 500) {
      return {
        bg: depth => getClass("bg", color, depth, defaultDepth),
        border: depth => getClass("border", color, depth, defaultDepth),
        txt: depth => getClass("text", color, depth, defaultDepth),
        caret: depth => getClass("caret", color, depth, defaultDepth)
      };
    }

    class ClassBuilder {
      constructor(classes, defaultClasses) {
        this.defaults =
          typeof classes === "function" ? classes(defaultClasses) : classes;

        this.classes = this.defaults;
      }

      flush() {
        this.classes = this.defaults;

        return this;
      }

      get() {
        return this.classes;
      }

      replace(classes, cond = true) {
        if (cond && classes) {
          this.classes = Object.keys(classes).reduce(
            (acc, from) => acc.replace(new RegExp(from, "g"), classes[from]),
            this.classes
          );
        }

        return this;
      }

      remove(classes, cond = true) {
        if (cond && classes) {
          this.classes = classes
            .split(" ")
            .reduce(
              (acc, cur) => acc.replace(new RegExp(cur, "g"), ""),
              this.classes
            );
        }

        return this;
      }

      add(className, cond = true, defaultValue) {
        if (!cond || !className) return this;

        switch (typeof className) {
          case "string":
          default:
            this.classes += ` ${className} `;
            return this;
          case "function":
            this.classes += ` ${className(defaultValue)} `;
            return this;
        }
      }
    }

    function filterProps(reserved, props) {

      return Object.keys(props).reduce(
        (acc, cur) =>
          cur.includes("$$") || cur.includes("Class") || reserved.includes(cur)
            ? acc
            : { ...acc, [cur]: props[cur] },
        {}
      );
    }

    // Thanks Lagden! https://svelte.dev/repl/61d9178d2b9944f2aa2bfe31612ab09f?version=3.6.7
    function ripple(color, centered) {
      return function(event) {
        const target = event.currentTarget;
        const circle = document.createElement("span");
        const d = Math.max(target.clientWidth, target.clientHeight);

        const removeCircle = () => {
          circle.remove();
          circle.removeEventListener("animationend", removeCircle);
        };

        circle.addEventListener("animationend", removeCircle);
        circle.style.width = circle.style.height = `${d}px`;
        const rect = target.getBoundingClientRect();

        if (centered) {
          circle.classList.add(
            "absolute",
            "top-0",
            "left-0",
            "ripple-centered",
            `bg-${color}-transDark`
          );
        } else {
          circle.style.left = `${event.clientX - rect.left - d / 2}px`;
          circle.style.top = `${event.clientY - rect.top - d / 2}px`;

          circle.classList.add("ripple-normal", `bg-${color}-trans`);
        }

        circle.classList.add("ripple");

        target.appendChild(circle);
      };
    }

    function r(color = "primary", centered = false) {
      return function(node) {
        node.addEventListener("click", ripple(color, centered));

        return {
          onDestroy: () => node.removeEventListener("click")
        };
      };
    }

    /* node_modules/smelte/src/components/Button/Button.svelte generated by Svelte v3.14.0 */
    const file$4 = "node_modules/smelte/src/components/Button/Button.svelte";

    // (148:0) {:else}
    function create_else_block(ctx) {
    	let button;
    	let t;
    	let ripple_action;
    	let current;
    	let dispose;
    	let if_block = ctx.icon && create_if_block_2(ctx);
    	const default_slot_template = ctx.$$slots.default;
    	const default_slot = create_slot(default_slot_template, ctx, null);

    	let button_levels = [
    		{
    			class: "" + (ctx.classes + " " + ctx.className + " button")
    		},
    		ctx.props,
    		{ disabled: ctx.disabled }
    	];

    	let button_data = {};

    	for (let i = 0; i < button_levels.length; i += 1) {
    		button_data = assign(button_data, button_levels[i]);
    	}

    	const block_1 = {
    		c: function create() {
    			button = element("button");
    			if (if_block) if_block.c();
    			t = space();
    			if (default_slot) default_slot.c();
    			set_attributes(button, button_data);
    			toggle_class(button, "border-solid", ctx.outlined);
    			toggle_class(button, "rounded-full", ctx.icon);
    			toggle_class(button, "w-full", ctx.block);
    			toggle_class(button, "rounded", ctx.basic || ctx.outlined || ctx.text);
    			toggle_class(button, "button", !ctx.icon);
    			add_location(button, file$4, 148, 2, 3947);

    			dispose = [
    				listen_dev(button, "click", ctx.click_handler_3, false, false, false),
    				listen_dev(button, "click", ctx.click_handler_1, false, false, false),
    				listen_dev(button, "mouseover", ctx.mouseover_handler_1, false, false, false),
    				listen_dev(button, "*", ctx._handler_1, false, false, false)
    			];
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, button, anchor);
    			if (if_block) if_block.m(button, null);
    			append_dev(button, t);

    			if (default_slot) {
    				default_slot.m(button, null);
    			}

    			ripple_action = ctx.ripple.call(null, button) || ({});
    			current = true;
    		},
    		p: function update(changed, ctx) {
    			if (ctx.icon) {
    				if (if_block) {
    					if_block.p(changed, ctx);
    					transition_in(if_block, 1);
    				} else {
    					if_block = create_if_block_2(ctx);
    					if_block.c();
    					transition_in(if_block, 1);
    					if_block.m(button, t);
    				}
    			} else if (if_block) {
    				group_outros();

    				transition_out(if_block, 1, 1, () => {
    					if_block = null;
    				});

    				check_outros();
    			}

    			if (default_slot && default_slot.p && changed.$$scope) {
    				default_slot.p(get_slot_changes(default_slot_template, ctx, changed, null), get_slot_context(default_slot_template, ctx, null));
    			}

    			set_attributes(button, get_spread_update(button_levels, [
    				(changed.classes || changed.className) && ({
    					class: "" + (ctx.classes + " " + ctx.className + " button")
    				}),
    				changed.props && ctx.props,
    				changed.disabled && ({ disabled: ctx.disabled })
    			]));

    			toggle_class(button, "border-solid", ctx.outlined);
    			toggle_class(button, "rounded-full", ctx.icon);
    			toggle_class(button, "w-full", ctx.block);
    			toggle_class(button, "rounded", ctx.basic || ctx.outlined || ctx.text);
    			toggle_class(button, "button", !ctx.icon);
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(if_block);
    			transition_in(default_slot, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(if_block);
    			transition_out(default_slot, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(button);
    			if (if_block) if_block.d();
    			if (default_slot) default_slot.d(detaching);
    			if (ripple_action && is_function(ripple_action.destroy)) ripple_action.destroy();
    			run_all(dispose);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block: block_1,
    		id: create_else_block.name,
    		type: "else",
    		source: "(148:0) {:else}",
    		ctx
    	});

    	return block_1;
    }

    // (122:0) {#if href}
    function create_if_block$1(ctx) {
    	let a;
    	let button;
    	let t;
    	let ripple_action;
    	let current;
    	let dispose;
    	let if_block = ctx.icon && create_if_block_1(ctx);
    	const default_slot_template = ctx.$$slots.default;
    	const default_slot = create_slot(default_slot_template, ctx, null);

    	let button_levels = [
    		{
    			class: "" + (ctx.classes + " " + ctx.className + " button")
    		},
    		ctx.props,
    		{ disabled: ctx.disabled }
    	];

    	let button_data = {};

    	for (let i = 0; i < button_levels.length; i += 1) {
    		button_data = assign(button_data, button_levels[i]);
    	}

    	let a_levels = [{ href: ctx.href }, ctx.props];
    	let a_data = {};

    	for (let i = 0; i < a_levels.length; i += 1) {
    		a_data = assign(a_data, a_levels[i]);
    	}

    	const block_1 = {
    		c: function create() {
    			a = element("a");
    			button = element("button");
    			if (if_block) if_block.c();
    			t = space();
    			if (default_slot) default_slot.c();
    			set_attributes(button, button_data);
    			toggle_class(button, "border-solid", ctx.outlined);
    			toggle_class(button, "rounded-full", ctx.icon);
    			toggle_class(button, "w-full", ctx.block);
    			toggle_class(button, "rounded", ctx.basic || ctx.outlined || ctx.text);
    			toggle_class(button, "button", !ctx.icon);
    			add_location(button, file$4, 126, 4, 3456);
    			set_attributes(a, a_data);
    			add_location(a, file$4, 122, 2, 3419);

    			dispose = [
    				listen_dev(button, "click", ctx.click_handler_2, false, false, false),
    				listen_dev(button, "click", ctx.click_handler, false, false, false),
    				listen_dev(button, "mouseover", ctx.mouseover_handler, false, false, false),
    				listen_dev(button, "*", ctx._handler, false, false, false)
    			];
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, a, anchor);
    			append_dev(a, button);
    			if (if_block) if_block.m(button, null);
    			append_dev(button, t);

    			if (default_slot) {
    				default_slot.m(button, null);
    			}

    			ripple_action = ctx.ripple.call(null, button) || ({});
    			current = true;
    		},
    		p: function update(changed, ctx) {
    			if (ctx.icon) {
    				if (if_block) {
    					if_block.p(changed, ctx);
    					transition_in(if_block, 1);
    				} else {
    					if_block = create_if_block_1(ctx);
    					if_block.c();
    					transition_in(if_block, 1);
    					if_block.m(button, t);
    				}
    			} else if (if_block) {
    				group_outros();

    				transition_out(if_block, 1, 1, () => {
    					if_block = null;
    				});

    				check_outros();
    			}

    			if (default_slot && default_slot.p && changed.$$scope) {
    				default_slot.p(get_slot_changes(default_slot_template, ctx, changed, null), get_slot_context(default_slot_template, ctx, null));
    			}

    			set_attributes(button, get_spread_update(button_levels, [
    				(changed.classes || changed.className) && ({
    					class: "" + (ctx.classes + " " + ctx.className + " button")
    				}),
    				changed.props && ctx.props,
    				changed.disabled && ({ disabled: ctx.disabled })
    			]));

    			toggle_class(button, "border-solid", ctx.outlined);
    			toggle_class(button, "rounded-full", ctx.icon);
    			toggle_class(button, "w-full", ctx.block);
    			toggle_class(button, "rounded", ctx.basic || ctx.outlined || ctx.text);
    			toggle_class(button, "button", !ctx.icon);
    			set_attributes(a, get_spread_update(a_levels, [changed.href && ({ href: ctx.href }), changed.props && ctx.props]));
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(if_block);
    			transition_in(default_slot, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(if_block);
    			transition_out(default_slot, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(a);
    			if (if_block) if_block.d();
    			if (default_slot) default_slot.d(detaching);
    			if (ripple_action && is_function(ripple_action.destroy)) ripple_action.destroy();
    			run_all(dispose);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block: block_1,
    		id: create_if_block$1.name,
    		type: "if",
    		source: "(122:0) {#if href}",
    		ctx
    	});

    	return block_1;
    }

    // (164:4) {#if icon}
    function create_if_block_2(ctx) {
    	let current;

    	const icon_1 = new Icon({
    			props: {
    				class: ctx.iClasses,
    				small: ctx.small,
    				$$slots: { default: [create_default_slot_1] },
    				$$scope: { ctx }
    			},
    			$$inline: true
    		});

    	const block_1 = {
    		c: function create() {
    			create_component(icon_1.$$.fragment);
    		},
    		m: function mount(target, anchor) {
    			mount_component(icon_1, target, anchor);
    			current = true;
    		},
    		p: function update(changed, ctx) {
    			const icon_1_changes = {};
    			if (changed.iClasses) icon_1_changes.class = ctx.iClasses;
    			if (changed.small) icon_1_changes.small = ctx.small;

    			if (changed.$$scope || changed.icon) {
    				icon_1_changes.$$scope = { changed, ctx };
    			}

    			icon_1.$set(icon_1_changes);
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(icon_1.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(icon_1.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			destroy_component(icon_1, detaching);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block: block_1,
    		id: create_if_block_2.name,
    		type: "if",
    		source: "(164:4) {#if icon}",
    		ctx
    	});

    	return block_1;
    }

    // (165:6) <Icon class={iClasses} {small}>
    function create_default_slot_1(ctx) {
    	let t;

    	const block_1 = {
    		c: function create() {
    			t = text(ctx.icon);
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, t, anchor);
    		},
    		p: function update(changed, ctx) {
    			if (changed.icon) set_data_dev(t, ctx.icon);
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(t);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block: block_1,
    		id: create_default_slot_1.name,
    		type: "slot",
    		source: "(165:6) <Icon class={iClasses} {small}>",
    		ctx
    	});

    	return block_1;
    }

    // (142:6) {#if icon}
    function create_if_block_1(ctx) {
    	let current;

    	const icon_1 = new Icon({
    			props: {
    				class: ctx.iClasses,
    				small: ctx.small,
    				$$slots: { default: [create_default_slot] },
    				$$scope: { ctx }
    			},
    			$$inline: true
    		});

    	const block_1 = {
    		c: function create() {
    			create_component(icon_1.$$.fragment);
    		},
    		m: function mount(target, anchor) {
    			mount_component(icon_1, target, anchor);
    			current = true;
    		},
    		p: function update(changed, ctx) {
    			const icon_1_changes = {};
    			if (changed.iClasses) icon_1_changes.class = ctx.iClasses;
    			if (changed.small) icon_1_changes.small = ctx.small;

    			if (changed.$$scope || changed.icon) {
    				icon_1_changes.$$scope = { changed, ctx };
    			}

    			icon_1.$set(icon_1_changes);
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(icon_1.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(icon_1.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			destroy_component(icon_1, detaching);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block: block_1,
    		id: create_if_block_1.name,
    		type: "if",
    		source: "(142:6) {#if icon}",
    		ctx
    	});

    	return block_1;
    }

    // (143:8) <Icon class={iClasses} {small}>
    function create_default_slot(ctx) {
    	let t;

    	const block_1 = {
    		c: function create() {
    			t = text(ctx.icon);
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, t, anchor);
    		},
    		p: function update(changed, ctx) {
    			if (changed.icon) set_data_dev(t, ctx.icon);
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(t);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block: block_1,
    		id: create_default_slot.name,
    		type: "slot",
    		source: "(143:8) <Icon class={iClasses} {small}>",
    		ctx
    	});

    	return block_1;
    }

    function create_fragment$4(ctx) {
    	let current_block_type_index;
    	let if_block;
    	let if_block_anchor;
    	let current;
    	const if_block_creators = [create_if_block$1, create_else_block];
    	const if_blocks = [];

    	function select_block_type(changed, ctx) {
    		if (ctx.href) return 0;
    		return 1;
    	}

    	current_block_type_index = select_block_type(null, ctx);
    	if_block = if_blocks[current_block_type_index] = if_block_creators[current_block_type_index](ctx);

    	const block_1 = {
    		c: function create() {
    			if_block.c();
    			if_block_anchor = empty();
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			if_blocks[current_block_type_index].m(target, anchor);
    			insert_dev(target, if_block_anchor, anchor);
    			current = true;
    		},
    		p: function update(changed, ctx) {
    			let previous_block_index = current_block_type_index;
    			current_block_type_index = select_block_type(changed, ctx);

    			if (current_block_type_index === previous_block_index) {
    				if_blocks[current_block_type_index].p(changed, ctx);
    			} else {
    				group_outros();

    				transition_out(if_blocks[previous_block_index], 1, 1, () => {
    					if_blocks[previous_block_index] = null;
    				});

    				check_outros();
    				if_block = if_blocks[current_block_type_index];

    				if (!if_block) {
    					if_block = if_blocks[current_block_type_index] = if_block_creators[current_block_type_index](ctx);
    					if_block.c();
    				}

    				transition_in(if_block, 1);
    				if_block.m(if_block_anchor.parentNode, if_block_anchor);
    			}
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(if_block);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(if_block);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if_blocks[current_block_type_index].d(detaching);
    			if (detaching) detach_dev(if_block_anchor);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block: block_1,
    		id: create_fragment$4.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block_1;
    }

    let commonDefault = "py-2 px-4 uppercase text-sm font-medium relative overflow-hidden";
    let basicDefault = "text-white transition";
    let outlinedDefault = "bg-transparent border border-solid";
    let textDefault = "bg-transparent border-none px-3 hover:bg-transparent";
    let iconDefault = "p-4 m-4 flex items-center";
    let fabDefault = "px-4 hover:bg-transparent";
    let smallDefault = "p-1 h-4 w-4";
    let disabledDefault = "bg-gray-300 text-gray-500 elevation-none pointer-events-none hover:bg-gray-300 cursor-default";
    let elevationDefault = "hover:elevation-5 elevation-3";

    function instance$2($$self, $$props, $$invalidate) {
    	let { class: className = "" } = $$props;
    	let { value = false } = $$props;
    	let { outlined = false } = $$props;
    	let { text = false } = $$props;
    	let { block = false } = $$props;
    	let { disabled = false } = $$props;
    	let { icon = null } = $$props;
    	let { small = false } = $$props;
    	let { light = false } = $$props;
    	let { dark = false } = $$props;
    	let { flat = false } = $$props;
    	let { iconClass = "" } = $$props;
    	let { color = "primary" } = $$props;
    	let { href = null } = $$props;
    	let { remove = "" } = $$props;
    	let { add = "" } = $$props;
    	let { replace = {} } = $$props;
    	let { commonClasses = commonDefault } = $$props;
    	let { basicClasses = basicDefault } = $$props;
    	let { outlinedClasses = outlinedDefault } = $$props;
    	let { textClasses = textDefault } = $$props;
    	let { iconClasses = iconDefault } = $$props;
    	let { fabClasses = fabDefault } = $$props;
    	let { smallClasses = smallDefault } = $$props;
    	let { disabledClasses = disabledDefault } = $$props;
    	let { elevationClasses = elevationDefault } = $$props;
    	const fab = text && icon;
    	const basic = !outlined && !text && !fab;
    	const elevation = (basic || icon) && !disabled && !flat && !text;
    	let classes = "";
    	let iClasses = "";
    	let shade = 0;
    	const { bg, border, txt } = utils(color);
    	const cb = new ClassBuilder(commonClasses, commonDefault);
    	let iconCb;

    	if (icon) {
    		$$invalidate("iconCb", iconCb = new ClassBuilder(iconClass));
    	}

    	const ripple = r(text || fab || outlined ? color : "white");

    	const props = filterProps(
    		[
    			"outlined",
    			"text",
    			"color",
    			"block",
    			"disabled",
    			"icon",
    			"small",
    			"light",
    			"dark",
    			"flat"
    		],
    		$$props
    	);

    	let { $$slots = {}, $$scope } = $$props;

    	function click_handler(event) {
    		bubble($$self, event);
    	}

    	function mouseover_handler(event) {
    		bubble($$self, event);
    	}

    	function _handler(event) {
    		bubble($$self, event);
    	}

    	function click_handler_1(event) {
    		bubble($$self, event);
    	}

    	function mouseover_handler_1(event) {
    		bubble($$self, event);
    	}

    	function _handler_1(event) {
    		bubble($$self, event);
    	}

    	const click_handler_2 = () => $$invalidate("value", value = !value);
    	const click_handler_3 = () => $$invalidate("value", value = !value);

    	$$self.$set = $$new_props => {
    		$$invalidate("$$props", $$props = assign(assign({}, $$props), $$new_props));
    		if ("class" in $$new_props) $$invalidate("className", className = $$new_props.class);
    		if ("value" in $$new_props) $$invalidate("value", value = $$new_props.value);
    		if ("outlined" in $$new_props) $$invalidate("outlined", outlined = $$new_props.outlined);
    		if ("text" in $$new_props) $$invalidate("text", text = $$new_props.text);
    		if ("block" in $$new_props) $$invalidate("block", block = $$new_props.block);
    		if ("disabled" in $$new_props) $$invalidate("disabled", disabled = $$new_props.disabled);
    		if ("icon" in $$new_props) $$invalidate("icon", icon = $$new_props.icon);
    		if ("small" in $$new_props) $$invalidate("small", small = $$new_props.small);
    		if ("light" in $$new_props) $$invalidate("light", light = $$new_props.light);
    		if ("dark" in $$new_props) $$invalidate("dark", dark = $$new_props.dark);
    		if ("flat" in $$new_props) $$invalidate("flat", flat = $$new_props.flat);
    		if ("iconClass" in $$new_props) $$invalidate("iconClass", iconClass = $$new_props.iconClass);
    		if ("color" in $$new_props) $$invalidate("color", color = $$new_props.color);
    		if ("href" in $$new_props) $$invalidate("href", href = $$new_props.href);
    		if ("remove" in $$new_props) $$invalidate("remove", remove = $$new_props.remove);
    		if ("add" in $$new_props) $$invalidate("add", add = $$new_props.add);
    		if ("replace" in $$new_props) $$invalidate("replace", replace = $$new_props.replace);
    		if ("commonClasses" in $$new_props) $$invalidate("commonClasses", commonClasses = $$new_props.commonClasses);
    		if ("basicClasses" in $$new_props) $$invalidate("basicClasses", basicClasses = $$new_props.basicClasses);
    		if ("outlinedClasses" in $$new_props) $$invalidate("outlinedClasses", outlinedClasses = $$new_props.outlinedClasses);
    		if ("textClasses" in $$new_props) $$invalidate("textClasses", textClasses = $$new_props.textClasses);
    		if ("iconClasses" in $$new_props) $$invalidate("iconClasses", iconClasses = $$new_props.iconClasses);
    		if ("fabClasses" in $$new_props) $$invalidate("fabClasses", fabClasses = $$new_props.fabClasses);
    		if ("smallClasses" in $$new_props) $$invalidate("smallClasses", smallClasses = $$new_props.smallClasses);
    		if ("disabledClasses" in $$new_props) $$invalidate("disabledClasses", disabledClasses = $$new_props.disabledClasses);
    		if ("elevationClasses" in $$new_props) $$invalidate("elevationClasses", elevationClasses = $$new_props.elevationClasses);
    		if ("$$scope" in $$new_props) $$invalidate("$$scope", $$scope = $$new_props.$$scope);
    	};

    	$$self.$capture_state = () => {
    		return {
    			className,
    			value,
    			outlined,
    			text,
    			block,
    			disabled,
    			icon,
    			small,
    			light,
    			dark,
    			flat,
    			iconClass,
    			color,
    			href,
    			remove,
    			add,
    			replace,
    			commonDefault,
    			basicDefault,
    			outlinedDefault,
    			textDefault,
    			iconDefault,
    			fabDefault,
    			smallDefault,
    			disabledDefault,
    			elevationDefault,
    			commonClasses,
    			basicClasses,
    			outlinedClasses,
    			textClasses,
    			iconClasses,
    			fabClasses,
    			smallClasses,
    			disabledClasses,
    			elevationClasses,
    			classes,
    			iClasses,
    			shade,
    			iconCb,
    			normal,
    			lighter
    		};
    	};

    	$$self.$inject_state = $$new_props => {
    		$$invalidate("$$props", $$props = assign(assign({}, $$props), $$new_props));
    		if ("className" in $$props) $$invalidate("className", className = $$new_props.className);
    		if ("value" in $$props) $$invalidate("value", value = $$new_props.value);
    		if ("outlined" in $$props) $$invalidate("outlined", outlined = $$new_props.outlined);
    		if ("text" in $$props) $$invalidate("text", text = $$new_props.text);
    		if ("block" in $$props) $$invalidate("block", block = $$new_props.block);
    		if ("disabled" in $$props) $$invalidate("disabled", disabled = $$new_props.disabled);
    		if ("icon" in $$props) $$invalidate("icon", icon = $$new_props.icon);
    		if ("small" in $$props) $$invalidate("small", small = $$new_props.small);
    		if ("light" in $$props) $$invalidate("light", light = $$new_props.light);
    		if ("dark" in $$props) $$invalidate("dark", dark = $$new_props.dark);
    		if ("flat" in $$props) $$invalidate("flat", flat = $$new_props.flat);
    		if ("iconClass" in $$props) $$invalidate("iconClass", iconClass = $$new_props.iconClass);
    		if ("color" in $$props) $$invalidate("color", color = $$new_props.color);
    		if ("href" in $$props) $$invalidate("href", href = $$new_props.href);
    		if ("remove" in $$props) $$invalidate("remove", remove = $$new_props.remove);
    		if ("add" in $$props) $$invalidate("add", add = $$new_props.add);
    		if ("replace" in $$props) $$invalidate("replace", replace = $$new_props.replace);
    		if ("commonDefault" in $$props) commonDefault = $$new_props.commonDefault;
    		if ("basicDefault" in $$props) $$invalidate("basicDefault", basicDefault = $$new_props.basicDefault);
    		if ("outlinedDefault" in $$props) $$invalidate("outlinedDefault", outlinedDefault = $$new_props.outlinedDefault);
    		if ("textDefault" in $$props) $$invalidate("textDefault", textDefault = $$new_props.textDefault);
    		if ("iconDefault" in $$props) $$invalidate("iconDefault", iconDefault = $$new_props.iconDefault);
    		if ("fabDefault" in $$props) $$invalidate("fabDefault", fabDefault = $$new_props.fabDefault);
    		if ("smallDefault" in $$props) $$invalidate("smallDefault", smallDefault = $$new_props.smallDefault);
    		if ("disabledDefault" in $$props) $$invalidate("disabledDefault", disabledDefault = $$new_props.disabledDefault);
    		if ("elevationDefault" in $$props) $$invalidate("elevationDefault", elevationDefault = $$new_props.elevationDefault);
    		if ("commonClasses" in $$props) $$invalidate("commonClasses", commonClasses = $$new_props.commonClasses);
    		if ("basicClasses" in $$props) $$invalidate("basicClasses", basicClasses = $$new_props.basicClasses);
    		if ("outlinedClasses" in $$props) $$invalidate("outlinedClasses", outlinedClasses = $$new_props.outlinedClasses);
    		if ("textClasses" in $$props) $$invalidate("textClasses", textClasses = $$new_props.textClasses);
    		if ("iconClasses" in $$props) $$invalidate("iconClasses", iconClasses = $$new_props.iconClasses);
    		if ("fabClasses" in $$props) $$invalidate("fabClasses", fabClasses = $$new_props.fabClasses);
    		if ("smallClasses" in $$props) $$invalidate("smallClasses", smallClasses = $$new_props.smallClasses);
    		if ("disabledClasses" in $$props) $$invalidate("disabledClasses", disabledClasses = $$new_props.disabledClasses);
    		if ("elevationClasses" in $$props) $$invalidate("elevationClasses", elevationClasses = $$new_props.elevationClasses);
    		if ("classes" in $$props) $$invalidate("classes", classes = $$new_props.classes);
    		if ("iClasses" in $$props) $$invalidate("iClasses", iClasses = $$new_props.iClasses);
    		if ("shade" in $$props) $$invalidate("shade", shade = $$new_props.shade);
    		if ("iconCb" in $$props) $$invalidate("iconCb", iconCb = $$new_props.iconCb);
    		if ("normal" in $$props) $$invalidate("normal", normal = $$new_props.normal);
    		if ("lighter" in $$props) $$invalidate("lighter", lighter = $$new_props.lighter);
    	};

    	let normal;
    	let lighter;

    	$$self.$$.update = (changed = { light: 1, dark: 1, shade: 1, basicClasses: 1, basicDefault: 1, normal: 1, lighter: 1, elevationClasses: 1, elevationDefault: 1, outlinedClasses: 1, outlined: 1, outlinedDefault: 1, text: 1, textClasses: 1, textDefault: 1, iconClasses: 1, icon: 1, iconDefault: 1, fabClasses: 1, fabDefault: 1, disabledClasses: 1, disabled: 1, disabledDefault: 1, smallClasses: 1, small: 1, smallDefault: 1, remove: 1, replace: 1, add: 1, iconCb: 1, iconClass: 1 }) => {
    		if (changed.light || changed.dark || changed.shade) {
    			 {
    				$$invalidate("shade", shade = light ? 200 : 0);
    				$$invalidate("shade", shade = dark ? -400 : shade);
    			}
    		}

    		if (changed.shade) {
    			 $$invalidate("normal", normal = 500 - shade);
    		}

    		if (changed.shade) {
    			 $$invalidate("lighter", lighter = 400 - shade);
    		}

    		if (changed.basicClasses || changed.basicDefault || changed.normal || changed.lighter || changed.elevationClasses || changed.elevationDefault || changed.outlinedClasses || changed.outlined || changed.outlinedDefault || changed.text || changed.textClasses || changed.textDefault || changed.iconClasses || changed.icon || changed.iconDefault || changed.fabClasses || changed.fabDefault || changed.disabledClasses || changed.disabled || changed.disabledDefault || changed.smallClasses || changed.small || changed.smallDefault || changed.remove || changed.replace || changed.add) {
    			 {
    				$$invalidate("classes", classes = cb.flush().add(basicClasses, basic, basicDefault).add(`${bg(normal)} hover:${bg(lighter)}`, basic).add(elevationClasses, elevation, elevationDefault).add(outlinedClasses, outlined, outlinedDefault).add(`${border(lighter)} ${txt(normal)} hover:${bg("trans")}`, outlined).add(`${txt(lighter)}`, text).add(textClasses, text, textDefault).add(iconClasses, icon, iconDefault).remove("py-2", icon).add(fabClasses, fab, fabDefault).remove(txt(lighter), fab).add(disabledClasses, disabled, disabledDefault).add(smallClasses, small, smallDefault).add("flex items-center justify-center", small && icon).remove(remove).replace(replace).add(add).get());
    			}
    		}

    		if (changed.iconCb || changed.iconClass) {
    			 {
    				if (iconCb) {
    					$$invalidate("iClasses", iClasses = iconCb.flush().add(txt(), fab && !iconClass).get());
    				}
    			}
    		}
    	};

    	return {
    		className,
    		value,
    		outlined,
    		text,
    		block,
    		disabled,
    		icon,
    		small,
    		light,
    		dark,
    		flat,
    		iconClass,
    		color,
    		href,
    		remove,
    		add,
    		replace,
    		commonClasses,
    		basicClasses,
    		outlinedClasses,
    		textClasses,
    		iconClasses,
    		fabClasses,
    		smallClasses,
    		disabledClasses,
    		elevationClasses,
    		basic,
    		classes,
    		iClasses,
    		ripple,
    		props,
    		click_handler,
    		mouseover_handler,
    		_handler,
    		click_handler_1,
    		mouseover_handler_1,
    		_handler_1,
    		click_handler_2,
    		click_handler_3,
    		$$props: $$props = exclude_internal_props($$props),
    		$$slots,
    		$$scope
    	};
    }

    class Button extends SvelteComponentDev {
    	constructor(options) {
    		super(options);

    		init(this, options, instance$2, create_fragment$4, safe_not_equal, {
    			class: "className",
    			value: 0,
    			outlined: 0,
    			text: 0,
    			block: 0,
    			disabled: 0,
    			icon: 0,
    			small: 0,
    			light: 0,
    			dark: 0,
    			flat: 0,
    			iconClass: 0,
    			color: 0,
    			href: 0,
    			remove: 0,
    			add: 0,
    			replace: 0,
    			commonClasses: 0,
    			basicClasses: 0,
    			outlinedClasses: 0,
    			textClasses: 0,
    			iconClasses: 0,
    			fabClasses: 0,
    			smallClasses: 0,
    			disabledClasses: 0,
    			elevationClasses: 0
    		});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "Button",
    			options,
    			id: create_fragment$4.name
    		});
    	}

    	get class() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set class(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get value() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set value(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get outlined() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set outlined(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get text() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set text(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get block() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set block(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get disabled() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set disabled(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get icon() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set icon(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get small() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set small(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get light() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set light(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get dark() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set dark(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get flat() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set flat(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get iconClass() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set iconClass(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get color() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set color(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get href() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set href(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get remove() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set remove(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get add() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set add(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get replace() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set replace(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get commonClasses() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set commonClasses(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get basicClasses() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set basicClasses(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get outlinedClasses() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set outlinedClasses(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get textClasses() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set textClasses(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get iconClasses() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set iconClasses(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get fabClasses() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set fabClasses(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get smallClasses() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set smallClasses(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get disabledClasses() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set disabledClasses(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get elevationClasses() {
    		throw new Error("<Button>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set elevationClasses(value) {
    		throw new Error("<Button>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}
    }

    /* node_modules/smelte/src/components/Util/Spacer.svelte generated by Svelte v3.14.0 */

    const file$5 = "node_modules/smelte/src/components/Util/Spacer.svelte";

    function create_fragment$5(ctx) {
    	let div;

    	const block = {
    		c: function create() {
    			div = element("div");
    			attr_dev(div, "class", "flex-grow");
    			add_location(div, file$5, 0, 0, 0);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div, anchor);
    		},
    		p: noop,
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$5.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class Spacer extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$5, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "Spacer",
    			options,
    			id: create_fragment$5.name
    		});
    	}
    }

    const Spacer$1 = Spacer;

    /* node_modules/smelte/src/components/List/ListItem.svelte generated by Svelte v3.14.0 */
    const file$6 = "node_modules/smelte/src/components/List/ListItem.svelte";

    // (50:2) {#if icon}
    function create_if_block_1$1(ctx) {
    	let current;

    	const icon_1 = new Icon({
    			props: {
    				class: "pr-6",
    				small: ctx.dense,
    				color: ctx.selected && ctx.navigation ? "text-primary-500" : "",
    				$$slots: { default: [create_default_slot$1] },
    				$$scope: { ctx }
    			},
    			$$inline: true
    		});

    	const block = {
    		c: function create() {
    			create_component(icon_1.$$.fragment);
    		},
    		m: function mount(target, anchor) {
    			mount_component(icon_1, target, anchor);
    			current = true;
    		},
    		p: function update(changed, ctx) {
    			const icon_1_changes = {};
    			if (changed.dense) icon_1_changes.small = ctx.dense;
    			if (changed.selected || changed.navigation) icon_1_changes.color = ctx.selected && ctx.navigation ? "text-primary-500" : "";

    			if (changed.$$scope || changed.icon) {
    				icon_1_changes.$$scope = { changed, ctx };
    			}

    			icon_1.$set(icon_1_changes);
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(icon_1.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(icon_1.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			destroy_component(icon_1, detaching);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_if_block_1$1.name,
    		type: "if",
    		source: "(50:2) {#if icon}",
    		ctx
    	});

    	return block;
    }

    // (51:4) <Icon       class="pr-6"       small={dense}       color={selected && navigation ? 'text-primary-500' : ''}>
    function create_default_slot$1(ctx) {
    	let t;

    	const block = {
    		c: function create() {
    			t = text(ctx.icon);
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, t, anchor);
    		},
    		p: function update(changed, ctx) {
    			if (changed.icon) set_data_dev(t, ctx.icon);
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(t);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_default_slot$1.name,
    		type: "slot",
    		source: "(51:4) <Icon       class=\\\"pr-6\\\"       small={dense}       color={selected && navigation ? 'text-primary-500' : ''}>",
    		ctx
    	});

    	return block;
    }

    // (63:4) {#if subheading}
    function create_if_block$2(ctx) {
    	let div;
    	let t;

    	const block = {
    		c: function create() {
    			div = element("div");
    			t = text(ctx.subheading);
    			attr_dev(div, "class", "text-gray-600 p-0 text-sm");
    			add_location(div, file$6, 63, 6, 1955);
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div, anchor);
    			append_dev(div, t);
    		},
    		p: function update(changed, ctx) {
    			if (changed.subheading) set_data_dev(t, ctx.subheading);
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_if_block$2.name,
    		type: "if",
    		source: "(63:4) {#if subheading}",
    		ctx
    	});

    	return block;
    }

    function create_fragment$6(ctx) {
    	let li;
    	let t0;
    	let div1;
    	let div0;
    	let t1;
    	let t2;
    	let li_class_value;
    	let ripple_action;
    	let current;
    	let dispose;
    	let if_block0 = ctx.icon && create_if_block_1$1(ctx);
    	const default_slot_template = ctx.$$slots.default;
    	const default_slot = create_slot(default_slot_template, ctx, null);
    	let if_block1 = ctx.subheading && create_if_block$2(ctx);

    	const block = {
    		c: function create() {
    			li = element("li");
    			if (if_block0) if_block0.c();
    			t0 = space();
    			div1 = element("div");
    			div0 = element("div");

    			if (!default_slot) {
    				t1 = text(ctx.text);
    			}

    			if (default_slot) default_slot.c();
    			t2 = space();
    			if (if_block1) if_block1.c();
    			attr_dev(div0, "class", ctx.itemClasses);
    			add_location(div0, file$6, 59, 4, 1865);
    			attr_dev(div1, "class", "flex flex-col p-0");
    			add_location(div1, file$6, 58, 2, 1829);
    			attr_dev(li, "class", li_class_value = "" + (ctx.basicClasses + " " + (ctx.selected ? ctx.selectedClasses : "") + " svelte-1pvx38n"));
    			attr_dev(li, "tabindex", ctx.tabindex);
    			toggle_class(li, "text-sm", ctx.navigation);
    			toggle_class(li, "py-2", ctx.dense);
    			toggle_class(li, "text-gray-600", ctx.disabled);
    			add_location(li, file$6, 39, 0, 1440);

    			dispose = [
    				listen_dev(li, "keypress", ctx.change, false, false, false),
    				listen_dev(li, "click", ctx.change, false, false, false),
    				listen_dev(li, "click", ctx.click_handler, false, false, false)
    			];
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, li, anchor);
    			if (if_block0) if_block0.m(li, null);
    			append_dev(li, t0);
    			append_dev(li, div1);
    			append_dev(div1, div0);

    			if (!default_slot) {
    				append_dev(div0, t1);
    			}

    			if (default_slot) {
    				default_slot.m(div0, null);
    			}

    			append_dev(div1, t2);
    			if (if_block1) if_block1.m(div1, null);
    			ripple_action = ctx.ripple.call(null, li) || ({});
    			current = true;
    		},
    		p: function update(changed, ctx) {
    			if (ctx.icon) {
    				if (if_block0) {
    					if_block0.p(changed, ctx);
    					transition_in(if_block0, 1);
    				} else {
    					if_block0 = create_if_block_1$1(ctx);
    					if_block0.c();
    					transition_in(if_block0, 1);
    					if_block0.m(li, t0);
    				}
    			} else if (if_block0) {
    				group_outros();

    				transition_out(if_block0, 1, 1, () => {
    					if_block0 = null;
    				});

    				check_outros();
    			}

    			if (!default_slot) {
    				if (!current || changed.text) set_data_dev(t1, ctx.text);
    			}

    			if (default_slot && default_slot.p && changed.$$scope) {
    				default_slot.p(get_slot_changes(default_slot_template, ctx, changed, null), get_slot_context(default_slot_template, ctx, null));
    			}

    			if (!current || changed.itemClasses) {
    				attr_dev(div0, "class", ctx.itemClasses);
    			}

    			if (ctx.subheading) {
    				if (if_block1) {
    					if_block1.p(changed, ctx);
    				} else {
    					if_block1 = create_if_block$2(ctx);
    					if_block1.c();
    					if_block1.m(div1, null);
    				}
    			} else if (if_block1) {
    				if_block1.d(1);
    				if_block1 = null;
    			}

    			if (!current || (changed.basicClasses || changed.selected || changed.selectedClasses) && li_class_value !== (li_class_value = "" + (ctx.basicClasses + " " + (ctx.selected ? ctx.selectedClasses : "") + " svelte-1pvx38n"))) {
    				attr_dev(li, "class", li_class_value);
    			}

    			if (!current || changed.tabindex) {
    				attr_dev(li, "tabindex", ctx.tabindex);
    			}

    			if (changed.basicClasses || changed.selected || changed.selectedClasses || changed.navigation) {
    				toggle_class(li, "text-sm", ctx.navigation);
    			}

    			if (changed.basicClasses || changed.selected || changed.selectedClasses || changed.dense) {
    				toggle_class(li, "py-2", ctx.dense);
    			}

    			if (changed.basicClasses || changed.selected || changed.selectedClasses || changed.disabled) {
    				toggle_class(li, "text-gray-600", ctx.disabled);
    			}
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(if_block0);
    			transition_in(default_slot, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(if_block0);
    			transition_out(default_slot, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(li);
    			if (if_block0) if_block0.d();
    			if (default_slot) default_slot.d(detaching);
    			if (if_block1) if_block1.d();
    			if (ripple_action && is_function(ripple_action.destroy)) ripple_action.destroy();
    			run_all(dispose);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$6.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    function instance$3($$self, $$props, $$invalidate) {
    	let { icon = "" } = $$props;
    	let { id = "" } = $$props;
    	let { value = "" } = $$props;
    	let { text = "" } = $$props;
    	let { subheading = "" } = $$props;
    	let { disabled = false } = $$props;
    	let { dense = false } = $$props;
    	let { navigation = false } = $$props;
    	let { to = "" } = $$props;
    	let { selected = false } = $$props;
    	let { tabindex = null } = $$props;
    	let { item = null } = $$props;
    	let { items = [] } = $$props;
    	let { level = null } = $$props;
    	let { basicClasses = "hover:bg-gray-transDark relative overflow-hidden transition p-4 cursor-pointer text-gray-700 flex items-center z-10" } = $$props;
    	let { itemClasses = "" } = $$props;
    	let { selectedClasses = "bg-gray-200 hover:bg-primary-transDark" } = $$props;
    	const ripple = r();
    	const dispatch = createEventDispatcher();

    	function change() {
    		if (disabled) return;
    		$$invalidate("value", value = id);
    		dispatch("change", id);
    	}

    	const writable_props = [
    		"icon",
    		"id",
    		"value",
    		"text",
    		"subheading",
    		"disabled",
    		"dense",
    		"navigation",
    		"to",
    		"selected",
    		"tabindex",
    		"item",
    		"items",
    		"level",
    		"basicClasses",
    		"itemClasses",
    		"selectedClasses"
    	];

    	Object.keys($$props).forEach(key => {
    		if (!writable_props.includes(key) && !key.startsWith("$$")) console.warn(`<ListItem> was created with unknown prop '${key}'`);
    	});

    	let { $$slots = {}, $$scope } = $$props;

    	function click_handler(event) {
    		bubble($$self, event);
    	}

    	$$self.$set = $$props => {
    		if ("icon" in $$props) $$invalidate("icon", icon = $$props.icon);
    		if ("id" in $$props) $$invalidate("id", id = $$props.id);
    		if ("value" in $$props) $$invalidate("value", value = $$props.value);
    		if ("text" in $$props) $$invalidate("text", text = $$props.text);
    		if ("subheading" in $$props) $$invalidate("subheading", subheading = $$props.subheading);
    		if ("disabled" in $$props) $$invalidate("disabled", disabled = $$props.disabled);
    		if ("dense" in $$props) $$invalidate("dense", dense = $$props.dense);
    		if ("navigation" in $$props) $$invalidate("navigation", navigation = $$props.navigation);
    		if ("to" in $$props) $$invalidate("to", to = $$props.to);
    		if ("selected" in $$props) $$invalidate("selected", selected = $$props.selected);
    		if ("tabindex" in $$props) $$invalidate("tabindex", tabindex = $$props.tabindex);
    		if ("item" in $$props) $$invalidate("item", item = $$props.item);
    		if ("items" in $$props) $$invalidate("items", items = $$props.items);
    		if ("level" in $$props) $$invalidate("level", level = $$props.level);
    		if ("basicClasses" in $$props) $$invalidate("basicClasses", basicClasses = $$props.basicClasses);
    		if ("itemClasses" in $$props) $$invalidate("itemClasses", itemClasses = $$props.itemClasses);
    		if ("selectedClasses" in $$props) $$invalidate("selectedClasses", selectedClasses = $$props.selectedClasses);
    		if ("$$scope" in $$props) $$invalidate("$$scope", $$scope = $$props.$$scope);
    	};

    	$$self.$capture_state = () => {
    		return {
    			icon,
    			id,
    			value,
    			text,
    			subheading,
    			disabled,
    			dense,
    			navigation,
    			to,
    			selected,
    			tabindex,
    			item,
    			items,
    			level,
    			basicClasses,
    			itemClasses,
    			selectedClasses
    		};
    	};

    	$$self.$inject_state = $$props => {
    		if ("icon" in $$props) $$invalidate("icon", icon = $$props.icon);
    		if ("id" in $$props) $$invalidate("id", id = $$props.id);
    		if ("value" in $$props) $$invalidate("value", value = $$props.value);
    		if ("text" in $$props) $$invalidate("text", text = $$props.text);
    		if ("subheading" in $$props) $$invalidate("subheading", subheading = $$props.subheading);
    		if ("disabled" in $$props) $$invalidate("disabled", disabled = $$props.disabled);
    		if ("dense" in $$props) $$invalidate("dense", dense = $$props.dense);
    		if ("navigation" in $$props) $$invalidate("navigation", navigation = $$props.navigation);
    		if ("to" in $$props) $$invalidate("to", to = $$props.to);
    		if ("selected" in $$props) $$invalidate("selected", selected = $$props.selected);
    		if ("tabindex" in $$props) $$invalidate("tabindex", tabindex = $$props.tabindex);
    		if ("item" in $$props) $$invalidate("item", item = $$props.item);
    		if ("items" in $$props) $$invalidate("items", items = $$props.items);
    		if ("level" in $$props) $$invalidate("level", level = $$props.level);
    		if ("basicClasses" in $$props) $$invalidate("basicClasses", basicClasses = $$props.basicClasses);
    		if ("itemClasses" in $$props) $$invalidate("itemClasses", itemClasses = $$props.itemClasses);
    		if ("selectedClasses" in $$props) $$invalidate("selectedClasses", selectedClasses = $$props.selectedClasses);
    	};

    	return {
    		icon,
    		id,
    		value,
    		text,
    		subheading,
    		disabled,
    		dense,
    		navigation,
    		to,
    		selected,
    		tabindex,
    		item,
    		items,
    		level,
    		basicClasses,
    		itemClasses,
    		selectedClasses,
    		ripple,
    		change,
    		click_handler,
    		$$slots,
    		$$scope
    	};
    }

    class ListItem extends SvelteComponentDev {
    	constructor(options) {
    		super(options);

    		init(this, options, instance$3, create_fragment$6, safe_not_equal, {
    			icon: 0,
    			id: 0,
    			value: 0,
    			text: 0,
    			subheading: 0,
    			disabled: 0,
    			dense: 0,
    			navigation: 0,
    			to: 0,
    			selected: 0,
    			tabindex: 0,
    			item: 0,
    			items: 0,
    			level: 0,
    			basicClasses: 0,
    			itemClasses: 0,
    			selectedClasses: 0
    		});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "ListItem",
    			options,
    			id: create_fragment$6.name
    		});
    	}

    	get icon() {
    		throw new Error("<ListItem>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set icon(value) {
    		throw new Error("<ListItem>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get id() {
    		throw new Error("<ListItem>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set id(value) {
    		throw new Error("<ListItem>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get value() {
    		throw new Error("<ListItem>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set value(value) {
    		throw new Error("<ListItem>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get text() {
    		throw new Error("<ListItem>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set text(value) {
    		throw new Error("<ListItem>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get subheading() {
    		throw new Error("<ListItem>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set subheading(value) {
    		throw new Error("<ListItem>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get disabled() {
    		throw new Error("<ListItem>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set disabled(value) {
    		throw new Error("<ListItem>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get dense() {
    		throw new Error("<ListItem>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set dense(value) {
    		throw new Error("<ListItem>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get navigation() {
    		throw new Error("<ListItem>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set navigation(value) {
    		throw new Error("<ListItem>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get to() {
    		throw new Error("<ListItem>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set to(value) {
    		throw new Error("<ListItem>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get selected() {
    		throw new Error("<ListItem>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set selected(value) {
    		throw new Error("<ListItem>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get tabindex() {
    		throw new Error("<ListItem>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set tabindex(value) {
    		throw new Error("<ListItem>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get item() {
    		throw new Error("<ListItem>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set item(value) {
    		throw new Error("<ListItem>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get items() {
    		throw new Error("<ListItem>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set items(value) {
    		throw new Error("<ListItem>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get level() {
    		throw new Error("<ListItem>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set level(value) {
    		throw new Error("<ListItem>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get basicClasses() {
    		throw new Error("<ListItem>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set basicClasses(value) {
    		throw new Error("<ListItem>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get itemClasses() {
    		throw new Error("<ListItem>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set itemClasses(value) {
    		throw new Error("<ListItem>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get selectedClasses() {
    		throw new Error("<ListItem>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set selectedClasses(value) {
    		throw new Error("<ListItem>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}
    }

    /* node_modules/smelte/src/components/List/List.svelte generated by Svelte v3.14.0 */
    const file$7 = "node_modules/smelte/src/components/List/List.svelte";
    const get_item_slot_changes_1 = ({ item, items, dense, navigation, value }) => ({ item: items, dense, navigation, value });
    const get_item_slot_context_1 = ({ item, items, dense, navigation, value }) => ({ item, dense, navigation, value });
    const get_item_slot_changes = ({ item, items, dense, navigation, value }) => ({ item: items, dense, navigation, value });
    const get_item_slot_context = ({ item, items, dense, navigation, value }) => ({ item, dense, navigation, value });

    function get_each_context$1(ctx, list, i) {
    	const child_ctx = Object.create(ctx);
    	child_ctx.item = list[i];
    	child_ctx.i = i;
    	return child_ctx;
    }

    // (42:6) {:else}
    function create_else_block$1(ctx) {
    	let updating_value;
    	let t;
    	let current;
    	const item_slot_template = ctx.$$slots.item;
    	const item_slot = create_slot(item_slot_template, ctx, get_item_slot_context_1);

    	const listitem_spread_levels = [
    		ctx.item,
    		{ tabindex: ctx.i + 1 },
    		{ id: ctx.id(ctx.item) },
    		{ selected: ctx.value === ctx.id(ctx.item) },
    		ctx.props
    	];

    	function listitem_value_binding_1(value_1) {
    		ctx.listitem_value_binding_1.call(null, value_1);
    	}

    	let listitem_props = {
    		$$slots: { default: [create_default_slot_1$1] },
    		$$scope: { ctx }
    	};

    	for (let i = 0; i < listitem_spread_levels.length; i += 1) {
    		listitem_props = assign(listitem_props, listitem_spread_levels[i]);
    	}

    	if (ctx.value !== void 0) {
    		listitem_props.value = ctx.value;
    	}

    	const listitem = new ListItem({ props: listitem_props, $$inline: true });
    	binding_callbacks.push(() => bind(listitem, "value", listitem_value_binding_1));
    	listitem.$on("change", ctx.change_handler_1);
    	listitem.$on("click", ctx.click_handler);

    	const block = {
    		c: function create() {
    			if (!item_slot) {
    				create_component(listitem.$$.fragment);
    				t = space();
    			}

    			if (item_slot) item_slot.c();
    		},
    		m: function mount(target, anchor) {
    			if (!item_slot) {
    				mount_component(listitem, target, anchor);
    				insert_dev(target, t, anchor);
    			}

    			if (item_slot) {
    				item_slot.m(target, anchor);
    			}

    			current = true;
    		},
    		p: function update(changed, ctx) {
    			if (!item_slot) {
    				const listitem_changes = changed.items || changed.id || changed.value || changed.props
    				? get_spread_update(listitem_spread_levels, [
    						changed.items && get_spread_object(ctx.item),
    						listitem_spread_levels[1],
    						(changed.id || changed.items) && ({ id: ctx.id(ctx.item) }),
    						(changed.value || changed.id || changed.items) && ({ selected: ctx.value === ctx.id(ctx.item) }),
    						changed.props && get_spread_object(ctx.props)
    					])
    				: {};

    				if (changed.$$scope || changed.items) {
    					listitem_changes.$$scope = { changed, ctx };
    				}

    				if (!updating_value && changed.value) {
    					updating_value = true;
    					listitem_changes.value = ctx.value;
    					add_flush_callback(() => updating_value = false);
    				}

    				listitem.$set(listitem_changes);
    			}

    			if (item_slot && item_slot.p && (changed.$$scope || changed.items || changed.dense || changed.navigation || changed.value)) {
    				item_slot.p(get_slot_changes(item_slot_template, ctx, changed, get_item_slot_changes_1), get_slot_context(item_slot_template, ctx, get_item_slot_context_1));
    			}
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(listitem.$$.fragment, local);
    			transition_in(item_slot, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(listitem.$$.fragment, local);
    			transition_out(item_slot, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (!item_slot) {
    				destroy_component(listitem, detaching);
    				if (detaching) detach_dev(t);
    			}

    			if (item_slot) item_slot.d(detaching);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_else_block$1.name,
    		type: "else",
    		source: "(42:6) {:else}",
    		ctx
    	});

    	return block;
    }

    // (34:6) {#if item.to}
    function create_if_block$3(ctx) {
    	let a;
    	let updating_value;
    	let a_tabindex_value;
    	let a_href_value;
    	let t;
    	let current;
    	const item_slot_template = ctx.$$slots.item;
    	const item_slot = create_slot(item_slot_template, ctx, get_item_slot_context);
    	const listitem_spread_levels = [ctx.item, { id: ctx.id(ctx.item) }, ctx.props];

    	function listitem_value_binding(value_1) {
    		ctx.listitem_value_binding.call(null, value_1);
    	}

    	let listitem_props = {
    		$$slots: { default: [create_default_slot$2] },
    		$$scope: { ctx }
    	};

    	for (let i = 0; i < listitem_spread_levels.length; i += 1) {
    		listitem_props = assign(listitem_props, listitem_spread_levels[i]);
    	}

    	if (ctx.value !== void 0) {
    		listitem_props.value = ctx.value;
    	}

    	const listitem = new ListItem({ props: listitem_props, $$inline: true });
    	binding_callbacks.push(() => bind(listitem, "value", listitem_value_binding));
    	listitem.$on("change", ctx.change_handler);

    	const block = {
    		c: function create() {
    			if (!item_slot) {
    				a = element("a");
    				create_component(listitem.$$.fragment);
    				t = space();
    			}

    			if (item_slot) item_slot.c();

    			if (!item_slot) {
    				attr_dev(a, "tabindex", a_tabindex_value = ctx.i + 1);
    				attr_dev(a, "href", a_href_value = ctx.item.to);
    				add_location(a, file$7, 35, 10, 860);
    			}
    		},
    		m: function mount(target, anchor) {
    			if (!item_slot) {
    				insert_dev(target, a, anchor);
    				mount_component(listitem, a, null);
    				insert_dev(target, t, anchor);
    			}

    			if (item_slot) {
    				item_slot.m(target, anchor);
    			}

    			current = true;
    		},
    		p: function update(changed, ctx) {
    			if (!item_slot) {
    				const listitem_changes = changed.items || changed.id || changed.props
    				? get_spread_update(listitem_spread_levels, [
    						changed.items && get_spread_object(ctx.item),
    						(changed.id || changed.items) && ({ id: ctx.id(ctx.item) }),
    						changed.props && get_spread_object(ctx.props)
    					])
    				: {};

    				if (changed.$$scope || changed.items) {
    					listitem_changes.$$scope = { changed, ctx };
    				}

    				if (!updating_value && changed.value) {
    					updating_value = true;
    					listitem_changes.value = ctx.value;
    					add_flush_callback(() => updating_value = false);
    				}

    				listitem.$set(listitem_changes);

    				if (!current || changed.items && a_href_value !== (a_href_value = ctx.item.to)) {
    					attr_dev(a, "href", a_href_value);
    				}
    			}

    			if (item_slot && item_slot.p && (changed.$$scope || changed.items || changed.dense || changed.navigation || changed.value)) {
    				item_slot.p(get_slot_changes(item_slot_template, ctx, changed, get_item_slot_changes), get_slot_context(item_slot_template, ctx, get_item_slot_context));
    			}
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(listitem.$$.fragment, local);
    			transition_in(item_slot, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(listitem.$$.fragment, local);
    			transition_out(item_slot, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (!item_slot) {
    				if (detaching) detach_dev(a);
    				destroy_component(listitem);
    				if (detaching) detach_dev(t);
    			}

    			if (item_slot) item_slot.d(detaching);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_if_block$3.name,
    		type: "if",
    		source: "(34:6) {#if item.to}",
    		ctx
    	});

    	return block;
    }

    // (44:10) <ListItem             bind:value             {...item}             tabindex={i + 1}             id={id(item)}             selected={value === id(item)}             {...props}             on:change             on:click>
    function create_default_slot_1$1(ctx) {
    	let t_value = getText(ctx.item) + "";
    	let t;

    	const block = {
    		c: function create() {
    			t = text(t_value);
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, t, anchor);
    		},
    		p: function update(changed, ctx) {
    			if (changed.items && t_value !== (t_value = getText(ctx.item) + "")) set_data_dev(t, t_value);
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(t);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_default_slot_1$1.name,
    		type: "slot",
    		source: "(44:10) <ListItem             bind:value             {...item}             tabindex={i + 1}             id={id(item)}             selected={value === id(item)}             {...props}             on:change             on:click>",
    		ctx
    	});

    	return block;
    }

    // (37:12) <ListItem bind:value {...item} id={id(item)} {...props} on:change>
    function create_default_slot$2(ctx) {
    	let t_value = ctx.item.text + "";
    	let t;

    	const block = {
    		c: function create() {
    			t = text(t_value);
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, t, anchor);
    		},
    		p: function update(changed, ctx) {
    			if (changed.items && t_value !== (t_value = ctx.item.text + "")) set_data_dev(t, t_value);
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(t);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_default_slot$2.name,
    		type: "slot",
    		source: "(37:12) <ListItem bind:value {...item} id={id(item)} {...props} on:change>",
    		ctx
    	});

    	return block;
    }

    // (33:4) {#each items as item, i}
    function create_each_block$1(ctx) {
    	let current_block_type_index;
    	let if_block;
    	let if_block_anchor;
    	let current;
    	const if_block_creators = [create_if_block$3, create_else_block$1];
    	const if_blocks = [];

    	function select_block_type(changed, ctx) {
    		if (ctx.item.to) return 0;
    		return 1;
    	}

    	current_block_type_index = select_block_type(null, ctx);
    	if_block = if_blocks[current_block_type_index] = if_block_creators[current_block_type_index](ctx);

    	const block = {
    		c: function create() {
    			if_block.c();
    			if_block_anchor = empty();
    		},
    		m: function mount(target, anchor) {
    			if_blocks[current_block_type_index].m(target, anchor);
    			insert_dev(target, if_block_anchor, anchor);
    			current = true;
    		},
    		p: function update(changed, ctx) {
    			let previous_block_index = current_block_type_index;
    			current_block_type_index = select_block_type(changed, ctx);

    			if (current_block_type_index === previous_block_index) {
    				if_blocks[current_block_type_index].p(changed, ctx);
    			} else {
    				group_outros();

    				transition_out(if_blocks[previous_block_index], 1, 1, () => {
    					if_blocks[previous_block_index] = null;
    				});

    				check_outros();
    				if_block = if_blocks[current_block_type_index];

    				if (!if_block) {
    					if_block = if_blocks[current_block_type_index] = if_block_creators[current_block_type_index](ctx);
    					if_block.c();
    				}

    				transition_in(if_block, 1);
    				if_block.m(if_block_anchor.parentNode, if_block_anchor);
    			}
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(if_block);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(if_block);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if_blocks[current_block_type_index].d(detaching);
    			if (detaching) detach_dev(if_block_anchor);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_each_block$1.name,
    		type: "each",
    		source: "(33:4) {#each items as item, i}",
    		ctx
    	});

    	return block;
    }

    function create_fragment$7(ctx) {
    	let div;
    	let ul;
    	let current;
    	let each_value = ctx.items;
    	let each_blocks = [];

    	for (let i = 0; i < each_value.length; i += 1) {
    		each_blocks[i] = create_each_block$1(get_each_context$1(ctx, each_value, i));
    	}

    	const out = i => transition_out(each_blocks[i], 1, 1, () => {
    		each_blocks[i] = null;
    	});

    	const block = {
    		c: function create() {
    			div = element("div");
    			ul = element("ul");

    			for (let i = 0; i < each_blocks.length; i += 1) {
    				each_blocks[i].c();
    			}

    			attr_dev(ul, "class", ctx.listClasses);
    			toggle_class(ul, "rounded-t-none", ctx.select);
    			add_location(ul, file$7, 31, 2, 683);
    			attr_dev(div, "class", ctx.className);
    			add_location(div, file$7, 30, 0, 657);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div, anchor);
    			append_dev(div, ul);

    			for (let i = 0; i < each_blocks.length; i += 1) {
    				each_blocks[i].m(ul, null);
    			}

    			current = true;
    		},
    		p: function update(changed, ctx) {
    			if (changed.items || changed.id || changed.props || changed.value || changed.$$scope || changed.dense || changed.navigation || changed.getText) {
    				each_value = ctx.items;
    				let i;

    				for (i = 0; i < each_value.length; i += 1) {
    					const child_ctx = get_each_context$1(ctx, each_value, i);

    					if (each_blocks[i]) {
    						each_blocks[i].p(changed, child_ctx);
    						transition_in(each_blocks[i], 1);
    					} else {
    						each_blocks[i] = create_each_block$1(child_ctx);
    						each_blocks[i].c();
    						transition_in(each_blocks[i], 1);
    						each_blocks[i].m(ul, null);
    					}
    				}

    				group_outros();

    				for (i = each_value.length; i < each_blocks.length; i += 1) {
    					out(i);
    				}

    				check_outros();
    			}

    			if (!current || changed.listClasses) {
    				attr_dev(ul, "class", ctx.listClasses);
    			}

    			if (changed.listClasses || changed.select) {
    				toggle_class(ul, "rounded-t-none", ctx.select);
    			}

    			if (!current || changed.className) {
    				attr_dev(div, "class", ctx.className);
    			}
    		},
    		i: function intro(local) {
    			if (current) return;

    			for (let i = 0; i < each_value.length; i += 1) {
    				transition_in(each_blocks[i]);
    			}

    			current = true;
    		},
    		o: function outro(local) {
    			each_blocks = each_blocks.filter(Boolean);

    			for (let i = 0; i < each_blocks.length; i += 1) {
    				transition_out(each_blocks[i]);
    			}

    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div);
    			destroy_each(each_blocks, detaching);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$7.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    function getText(item) {
    	if (item.text !== undefined) return item.text;
    	if (item.value !== undefined) return item.value;
    	return item;
    }

    function instance$4($$self, $$props, $$invalidate) {
    	let { items = [] } = $$props;
    	let { item = {} } = $$props;
    	let { value = "" } = $$props;
    	let { text = "" } = $$props;
    	let { dense = false } = $$props;
    	let { navigation = false } = $$props;
    	let { level = null } = $$props;
    	let { select = false } = $$props;
    	let { class: className = "" } = $$props;
    	let { listClasses = "py-2 rounded" } = $$props;
    	const props = { dense, navigation };
    	const id = item => item.id || item.value || item.to || item.text || item;

    	const writable_props = [
    		"items",
    		"item",
    		"value",
    		"text",
    		"dense",
    		"navigation",
    		"level",
    		"select",
    		"class",
    		"listClasses"
    	];

    	Object.keys($$props).forEach(key => {
    		if (!writable_props.includes(key) && !key.startsWith("$$")) console.warn(`<List> was created with unknown prop '${key}'`);
    	});

    	let { $$slots = {}, $$scope } = $$props;

    	function listitem_value_binding(value_1) {
    		value = value_1;
    		$$invalidate("value", value);
    	}

    	function change_handler(event) {
    		bubble($$self, event);
    	}

    	function listitem_value_binding_1(value_1) {
    		value = value_1;
    		$$invalidate("value", value);
    	}

    	function change_handler_1(event) {
    		bubble($$self, event);
    	}

    	function click_handler(event) {
    		bubble($$self, event);
    	}

    	$$self.$set = $$props => {
    		if ("items" in $$props) $$invalidate("items", items = $$props.items);
    		if ("item" in $$props) $$invalidate("item", item = $$props.item);
    		if ("value" in $$props) $$invalidate("value", value = $$props.value);
    		if ("text" in $$props) $$invalidate("text", text = $$props.text);
    		if ("dense" in $$props) $$invalidate("dense", dense = $$props.dense);
    		if ("navigation" in $$props) $$invalidate("navigation", navigation = $$props.navigation);
    		if ("level" in $$props) $$invalidate("level", level = $$props.level);
    		if ("select" in $$props) $$invalidate("select", select = $$props.select);
    		if ("class" in $$props) $$invalidate("className", className = $$props.class);
    		if ("listClasses" in $$props) $$invalidate("listClasses", listClasses = $$props.listClasses);
    		if ("$$scope" in $$props) $$invalidate("$$scope", $$scope = $$props.$$scope);
    	};

    	$$self.$capture_state = () => {
    		return {
    			items,
    			item,
    			value,
    			text,
    			dense,
    			navigation,
    			level,
    			select,
    			className,
    			listClasses
    		};
    	};

    	$$self.$inject_state = $$props => {
    		if ("items" in $$props) $$invalidate("items", items = $$props.items);
    		if ("item" in $$props) $$invalidate("item", item = $$props.item);
    		if ("value" in $$props) $$invalidate("value", value = $$props.value);
    		if ("text" in $$props) $$invalidate("text", text = $$props.text);
    		if ("dense" in $$props) $$invalidate("dense", dense = $$props.dense);
    		if ("navigation" in $$props) $$invalidate("navigation", navigation = $$props.navigation);
    		if ("level" in $$props) $$invalidate("level", level = $$props.level);
    		if ("select" in $$props) $$invalidate("select", select = $$props.select);
    		if ("className" in $$props) $$invalidate("className", className = $$props.className);
    		if ("listClasses" in $$props) $$invalidate("listClasses", listClasses = $$props.listClasses);
    	};

    	return {
    		items,
    		item,
    		value,
    		text,
    		dense,
    		navigation,
    		level,
    		select,
    		className,
    		listClasses,
    		props,
    		id,
    		listitem_value_binding,
    		change_handler,
    		listitem_value_binding_1,
    		change_handler_1,
    		click_handler,
    		$$slots,
    		$$scope
    	};
    }

    class List extends SvelteComponentDev {
    	constructor(options) {
    		super(options);

    		init(this, options, instance$4, create_fragment$7, safe_not_equal, {
    			items: 0,
    			item: 0,
    			value: 0,
    			text: 0,
    			dense: 0,
    			navigation: 0,
    			level: 0,
    			select: 0,
    			class: "className",
    			listClasses: 0
    		});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "List",
    			options,
    			id: create_fragment$7.name
    		});
    	}

    	get items() {
    		throw new Error("<List>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set items(value) {
    		throw new Error("<List>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get item() {
    		throw new Error("<List>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set item(value) {
    		throw new Error("<List>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get value() {
    		throw new Error("<List>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set value(value) {
    		throw new Error("<List>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get text() {
    		throw new Error("<List>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set text(value) {
    		throw new Error("<List>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get dense() {
    		throw new Error("<List>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set dense(value) {
    		throw new Error("<List>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get navigation() {
    		throw new Error("<List>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set navigation(value) {
    		throw new Error("<List>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get level() {
    		throw new Error("<List>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set level(value) {
    		throw new Error("<List>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get select() {
    		throw new Error("<List>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set select(value) {
    		throw new Error("<List>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get class() {
    		throw new Error("<List>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set class(value) {
    		throw new Error("<List>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get listClasses() {
    		throw new Error("<List>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set listClasses(value) {
    		throw new Error("<List>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}
    }

    /* node_modules/smelte/src/components/TextField/TextField.svelte generated by Svelte v3.14.0 */
    const file$8 = "node_modules/smelte/src/components/TextField/TextField.svelte";
    const get_prepend_slot_changes = () => ({});
    const get_prepend_slot_context = () => ({});
    const get_append_slot_changes = () => ({});
    const get_append_slot_context = () => ({});

    // (162:4) {#if append}
    function create_if_block_5(ctx) {
    	let div;
    	let div_class_value;
    	let current;

    	const icon = new Icon({
    			props: {
    				reverse: ctx.appendReverse,
    				class: "" + ((ctx.focused ? ctx.txt() : "") + " " + ctx.iconClass),
    				$$slots: { default: [create_default_slot_1$2] },
    				$$scope: { ctx }
    			},
    			$$inline: true
    		});

    	const block = {
    		c: function create() {
    			div = element("div");
    			create_component(icon.$$.fragment);
    			attr_dev(div, "class", div_class_value = "" + (null_to_empty(ctx.aClasses) + " svelte-8az9n8"));
    			add_location(div, file$8, 162, 6, 5290);
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div, anchor);
    			mount_component(icon, div, null);
    			current = true;
    		},
    		p: function update(changed, ctx) {
    			const icon_changes = {};
    			if (changed.appendReverse) icon_changes.reverse = ctx.appendReverse;
    			if (changed.focused || changed.iconClass) icon_changes.class = "" + ((ctx.focused ? ctx.txt() : "") + " " + ctx.iconClass);

    			if (changed.$$scope || changed.append) {
    				icon_changes.$$scope = { changed, ctx };
    			}

    			icon.$set(icon_changes);

    			if (!current || changed.aClasses && div_class_value !== (div_class_value = "" + (null_to_empty(ctx.aClasses) + " svelte-8az9n8"))) {
    				attr_dev(div, "class", div_class_value);
    			}
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(icon.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(icon.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div);
    			destroy_component(icon);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_if_block_5.name,
    		type: "if",
    		source: "(162:4) {#if append}",
    		ctx
    	});

    	return block;
    }

    // (164:8) <Icon           reverse={appendReverse}           class="{focused ? txt() : ''} {iconClass}"         >
    function create_default_slot_1$2(ctx) {
    	let t;

    	const block = {
    		c: function create() {
    			t = text(ctx.append);
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, t, anchor);
    		},
    		p: function update(changed, ctx) {
    			if (changed.append) set_data_dev(t, ctx.append);
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(t);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_default_slot_1$2.name,
    		type: "slot",
    		source: "(164:8) <Icon           reverse={appendReverse}           class=\\\"{focused ? txt() : ''} {iconClass}\\\"         >",
    		ctx
    	});

    	return block;
    }

    // (202:38) 
    function create_if_block_4(ctx) {
    	let input;
    	let input_class_value;
    	let dispose;

    	const block = {
    		c: function create() {
    			input = element("input");
    			input.readOnly = true;
    			attr_dev(input, "class", input_class_value = "" + (null_to_empty(ctx.iClasses) + " svelte-8az9n8"));
    			input.value = ctx.value;
    			add_location(input, file$8, 202, 6, 6238);

    			dispose = [
    				listen_dev(input, "click", ctx.toggleFocused, false, false, false),
    				listen_dev(input, "change", ctx.change_handler_2, false, false, false),
    				listen_dev(input, "input", ctx.input_handler_2, false, false, false),
    				listen_dev(input, "click", ctx.click_handler_2, false, false, false),
    				listen_dev(input, "blur", ctx.blur_handler_2, false, false, false),
    				listen_dev(input, "focus", ctx.focus_handler_2, false, false, false)
    			];
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, input, anchor);
    		},
    		p: function update(changed, ctx) {
    			if (changed.iClasses && input_class_value !== (input_class_value = "" + (null_to_empty(ctx.iClasses) + " svelte-8az9n8"))) {
    				attr_dev(input, "class", input_class_value);
    			}

    			if (changed.value) {
    				prop_dev(input, "value", ctx.value);
    			}
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(input);
    			run_all(dispose);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_if_block_4.name,
    		type: "if",
    		source: "(202:38) ",
    		ctx
    	});

    	return block;
    }

    // (187:34) 
    function create_if_block_3(ctx) {
    	let textarea_1;
    	let dispose;

    	let textarea_1_levels = [
    		{ rows: ctx.rows },
    		{ "aria-label": ctx.label },
    		{ class: ctx.iClasses },
    		ctx.props,
    		{
    			placeholder: !ctx.value ? ctx.placeholder : ""
    		}
    	];

    	let textarea_1_data = {};

    	for (let i = 0; i < textarea_1_levels.length; i += 1) {
    		textarea_1_data = assign(textarea_1_data, textarea_1_levels[i]);
    	}

    	const block = {
    		c: function create() {
    			textarea_1 = element("textarea");
    			set_attributes(textarea_1, textarea_1_data);
    			toggle_class(textarea_1, "svelte-8az9n8", true);
    			add_location(textarea_1, file$8, 187, 6, 5877);

    			dispose = [
    				listen_dev(textarea_1, "input", ctx.textarea_1_input_handler),
    				listen_dev(textarea_1, "change", ctx.change_handler_1, false, false, false),
    				listen_dev(textarea_1, "input", ctx.input_handler_1, false, false, false),
    				listen_dev(textarea_1, "click", ctx.click_handler_1, false, false, false),
    				listen_dev(textarea_1, "focus", ctx.focus_handler_1, false, false, false),
    				listen_dev(textarea_1, "blur", ctx.blur_handler_1, false, false, false),
    				listen_dev(textarea_1, "focus", ctx.toggleFocused, false, false, false),
    				listen_dev(textarea_1, "blur", ctx.toggleFocused, false, false, false)
    			];
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, textarea_1, anchor);
    			set_input_value(textarea_1, ctx.value);
    		},
    		p: function update(changed, ctx) {
    			set_attributes(textarea_1, get_spread_update(textarea_1_levels, [
    				changed.rows && ({ rows: ctx.rows }),
    				changed.label && ({ "aria-label": ctx.label }),
    				changed.iClasses && ({ class: ctx.iClasses }),
    				changed.props && ctx.props,
    				(changed.value || changed.placeholder) && ({
    					placeholder: !ctx.value ? ctx.placeholder : ""
    				})
    			]));

    			if (changed.value) {
    				set_input_value(textarea_1, ctx.value);
    			}

    			toggle_class(textarea_1, "svelte-8az9n8", true);
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(textarea_1);
    			run_all(dispose);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_if_block_3.name,
    		type: "if",
    		source: "(187:34) ",
    		ctx
    	});

    	return block;
    }

    // (173:4) {#if (!textarea && !select) || autocomplete}
    function create_if_block_2$1(ctx) {
    	let input;
    	let dispose;

    	let input_levels = [
    		{ "aria-label": ctx.label },
    		{ class: ctx.iClasses },
    		ctx.props,
    		{
    			placeholder: !ctx.value ? ctx.placeholder : ""
    		}
    	];

    	let input_data = {};

    	for (let i = 0; i < input_levels.length; i += 1) {
    		input_data = assign(input_data, input_levels[i]);
    	}

    	const block = {
    		c: function create() {
    			input = element("input");
    			set_attributes(input, input_data);
    			toggle_class(input, "svelte-8az9n8", true);
    			add_location(input, file$8, 173, 6, 5538);

    			dispose = [
    				listen_dev(input, "input", ctx.input_input_handler),
    				listen_dev(input, "focus", ctx.toggleFocused, false, false, false),
    				listen_dev(input, "blur", ctx.toggleFocused, false, false, false),
    				listen_dev(input, "blur", ctx.blur_handler, false, false, false),
    				listen_dev(input, "change", ctx.change_handler, false, false, false),
    				listen_dev(input, "input", ctx.input_handler, false, false, false),
    				listen_dev(input, "click", ctx.click_handler, false, false, false),
    				listen_dev(input, "focus", ctx.focus_handler, false, false, false)
    			];
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, input, anchor);
    			set_input_value(input, ctx.value);
    		},
    		p: function update(changed, ctx) {
    			set_attributes(input, get_spread_update(input_levels, [
    				changed.label && ({ "aria-label": ctx.label }),
    				changed.iClasses && ({ class: ctx.iClasses }),
    				changed.props && ctx.props,
    				(changed.value || changed.placeholder) && ({
    					placeholder: !ctx.value ? ctx.placeholder : ""
    				})
    			]));

    			if (changed.value && input.value !== ctx.value) {
    				set_input_value(input, ctx.value);
    			}

    			toggle_class(input, "svelte-8az9n8", true);
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(input);
    			run_all(dispose);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_if_block_2$1.name,
    		type: "if",
    		source: "(173:4) {#if (!textarea && !select) || autocomplete}",
    		ctx
    	});

    	return block;
    }

    // (219:4) {#if prepend}
    function create_if_block_1$2(ctx) {
    	let div;
    	let div_class_value;
    	let current;

    	const icon = new Icon({
    			props: {
    				reverse: ctx.prependReverse,
    				class: "" + ((ctx.focused ? ctx.txt() : "") + " " + ctx.iconClass),
    				$$slots: { default: [create_default_slot$3] },
    				$$scope: { ctx }
    			},
    			$$inline: true
    		});

    	const block = {
    		c: function create() {
    			div = element("div");
    			create_component(icon.$$.fragment);
    			attr_dev(div, "class", div_class_value = "" + (null_to_empty(ctx.pClasses) + " svelte-8az9n8"));
    			add_location(div, file$8, 219, 6, 6530);
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div, anchor);
    			mount_component(icon, div, null);
    			current = true;
    		},
    		p: function update(changed, ctx) {
    			const icon_changes = {};
    			if (changed.prependReverse) icon_changes.reverse = ctx.prependReverse;
    			if (changed.focused || changed.iconClass) icon_changes.class = "" + ((ctx.focused ? ctx.txt() : "") + " " + ctx.iconClass);

    			if (changed.$$scope || changed.prepend) {
    				icon_changes.$$scope = { changed, ctx };
    			}

    			icon.$set(icon_changes);

    			if (!current || changed.pClasses && div_class_value !== (div_class_value = "" + (null_to_empty(ctx.pClasses) + " svelte-8az9n8"))) {
    				attr_dev(div, "class", div_class_value);
    			}
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(icon.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(icon.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div);
    			destroy_component(icon);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_if_block_1$2.name,
    		type: "if",
    		source: "(219:4) {#if prepend}",
    		ctx
    	});

    	return block;
    }

    // (221:8) <Icon           reverse={prependReverse}           class="{focused ? txt() : ''} {iconClass}"         >
    function create_default_slot$3(ctx) {
    	let t;

    	const block = {
    		c: function create() {
    			t = text(ctx.prepend);
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, t, anchor);
    		},
    		p: function update(changed, ctx) {
    			if (changed.prepend) set_data_dev(t, ctx.prepend);
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(t);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_default_slot$3.name,
    		type: "slot",
    		source: "(221:8) <Icon           reverse={prependReverse}           class=\\\"{focused ? txt() : ''} {iconClass}\\\"         >",
    		ctx
    	});

    	return block;
    }

    // (241:2) {#if showHint}
    function create_if_block$4(ctx) {
    	let div;
    	let t;
    	let div_transition;
    	let current;

    	const block = {
    		c: function create() {
    			div = element("div");
    			t = text(ctx.showHint);
    			attr_dev(div, "class", "text-xs py-1 pl-4 absolute bottom-0 left-0");
    			toggle_class(div, "text-gray-600", ctx.hint);
    			toggle_class(div, "text-error-500", ctx.error);
    			add_location(div, file$8, 241, 4, 7080);
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div, anchor);
    			append_dev(div, t);
    			current = true;
    		},
    		p: function update(changed, ctx) {
    			if (!current || changed.showHint) set_data_dev(t, ctx.showHint);

    			if (changed.hint) {
    				toggle_class(div, "text-gray-600", ctx.hint);
    			}

    			if (changed.error) {
    				toggle_class(div, "text-error-500", ctx.error);
    			}
    		},
    		i: function intro(local) {
    			if (current) return;

    			add_render_callback(() => {
    				if (!div_transition) div_transition = create_bidirectional_transition(div, fly, { y: -10, duration: 100, easing: quadOut }, true);
    				div_transition.run(1);
    			});

    			current = true;
    		},
    		o: function outro(local) {
    			if (!div_transition) div_transition = create_bidirectional_transition(div, fly, { y: -10, duration: 100, easing: quadOut }, false);
    			div_transition.run(0);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div);
    			if (detaching && div_transition) div_transition.end();
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_if_block$4.name,
    		type: "if",
    		source: "(241:2) {#if showHint}",
    		ctx
    	});

    	return block;
    }

    function create_fragment$8(ctx) {
    	let div5;
    	let div4;
    	let label_1;
    	let t0;
    	let label_1_class_value;
    	let t1;
    	let div0;
    	let div0_class_value;
    	let t2;
    	let t3;
    	let t4;
    	let div1;
    	let div1_class_value;
    	let t5;
    	let t6;
    	let div3;
    	let div2;
    	let div2_class_value;
    	let t7;
    	let div5_class_value;
    	let current;
    	let dispose;
    	const append_slot_template = ctx.$$slots.append;
    	const append_slot = create_slot(append_slot_template, ctx, get_append_slot_context);
    	let if_block0 = ctx.append && create_if_block_5(ctx);

    	function select_block_type(changed, ctx) {
    		if (!ctx.textarea && !ctx.select || ctx.autocomplete) return create_if_block_2$1;
    		if (ctx.textarea && !ctx.select) return create_if_block_3;
    		if (ctx.select && !ctx.autocomplete) return create_if_block_4;
    	}

    	let current_block_type = select_block_type(null, ctx);
    	let if_block1 = current_block_type && current_block_type(ctx);
    	const prepend_slot_template = ctx.$$slots.prepend;
    	const prepend_slot = create_slot(prepend_slot_template, ctx, get_prepend_slot_context);
    	let if_block2 = ctx.prepend && create_if_block_1$2(ctx);
    	let if_block3 = ctx.showHint && create_if_block$4(ctx);

    	const block = {
    		c: function create() {
    			div5 = element("div");
    			div4 = element("div");
    			label_1 = element("label");
    			t0 = text(ctx.label);
    			t1 = space();
    			div0 = element("div");
    			if (append_slot) append_slot.c();
    			t2 = space();
    			if (if_block0) if_block0.c();
    			t3 = space();
    			if (if_block1) if_block1.c();
    			t4 = space();
    			div1 = element("div");
    			if (prepend_slot) prepend_slot.c();
    			t5 = space();
    			if (if_block2) if_block2.c();
    			t6 = space();
    			div3 = element("div");
    			div2 = element("div");
    			t7 = space();
    			if (if_block3) if_block3.c();
    			attr_dev(label_1, "class", label_1_class_value = "" + (null_to_empty(ctx.lClasses) + " svelte-8az9n8"));
    			add_location(label_1, file$8, 153, 4, 5146);
    			attr_dev(div0, "class", div0_class_value = "" + (null_to_empty(ctx.aClasses) + " svelte-8az9n8"));
    			add_location(div0, file$8, 157, 4, 5203);
    			attr_dev(div1, "class", div1_class_value = "" + (null_to_empty(ctx.pClasses) + " svelte-8az9n8"));
    			add_location(div1, file$8, 214, 4, 6441);
    			attr_dev(div2, "class", div2_class_value = "mx-auto w-0 " + (ctx.focused ? ctx.bg() : "") + " svelte-8az9n8");
    			set_style(div2, "height", "2px");
    			set_style(div2, "transition", "width .2s ease");
    			toggle_class(div2, "w-full", ctx.focused || ctx.error);
    			toggle_class(div2, "bg-error-500", ctx.error);
    			add_location(div2, file$8, 232, 6, 6849);
    			attr_dev(div3, "class", "line absolute bottom-0 left-0 w-full bg-gray-600 svelte-8az9n8");
    			toggle_class(div3, "hidden", ctx.noUnderline || ctx.outlined);
    			add_location(div3, file$8, 229, 4, 6729);
    			attr_dev(div4, "class", "relative");
    			toggle_class(div4, "text-error-500", ctx.error);
    			add_location(div4, file$8, 152, 2, 5090);
    			attr_dev(div5, "class", div5_class_value = "" + (null_to_empty(ctx.wClasses) + " svelte-8az9n8"));
    			add_location(div5, file$8, 151, 0, 5065);
    			dispose = listen_dev(window, "click", ctx.click_handler_3, false, false, false);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div5, anchor);
    			append_dev(div5, div4);
    			append_dev(div4, label_1);
    			append_dev(label_1, t0);
    			append_dev(div4, t1);
    			append_dev(div4, div0);

    			if (append_slot) {
    				append_slot.m(div0, null);
    			}

    			append_dev(div4, t2);
    			if (if_block0) if_block0.m(div4, null);
    			append_dev(div4, t3);
    			if (if_block1) if_block1.m(div4, null);
    			append_dev(div4, t4);
    			append_dev(div4, div1);

    			if (prepend_slot) {
    				prepend_slot.m(div1, null);
    			}

    			append_dev(div4, t5);
    			if (if_block2) if_block2.m(div4, null);
    			append_dev(div4, t6);
    			append_dev(div4, div3);
    			append_dev(div3, div2);
    			append_dev(div5, t7);
    			if (if_block3) if_block3.m(div5, null);
    			current = true;
    		},
    		p: function update(changed, ctx) {
    			if (!current || changed.label) set_data_dev(t0, ctx.label);

    			if (!current || changed.lClasses && label_1_class_value !== (label_1_class_value = "" + (null_to_empty(ctx.lClasses) + " svelte-8az9n8"))) {
    				attr_dev(label_1, "class", label_1_class_value);
    			}

    			if (append_slot && append_slot.p && changed.$$scope) {
    				append_slot.p(get_slot_changes(append_slot_template, ctx, changed, get_append_slot_changes), get_slot_context(append_slot_template, ctx, get_append_slot_context));
    			}

    			if (!current || changed.aClasses && div0_class_value !== (div0_class_value = "" + (null_to_empty(ctx.aClasses) + " svelte-8az9n8"))) {
    				attr_dev(div0, "class", div0_class_value);
    			}

    			if (ctx.append) {
    				if (if_block0) {
    					if_block0.p(changed, ctx);
    					transition_in(if_block0, 1);
    				} else {
    					if_block0 = create_if_block_5(ctx);
    					if_block0.c();
    					transition_in(if_block0, 1);
    					if_block0.m(div4, t3);
    				}
    			} else if (if_block0) {
    				group_outros();

    				transition_out(if_block0, 1, 1, () => {
    					if_block0 = null;
    				});

    				check_outros();
    			}

    			if (current_block_type === (current_block_type = select_block_type(changed, ctx)) && if_block1) {
    				if_block1.p(changed, ctx);
    			} else {
    				if (if_block1) if_block1.d(1);
    				if_block1 = current_block_type && current_block_type(ctx);

    				if (if_block1) {
    					if_block1.c();
    					if_block1.m(div4, t4);
    				}
    			}

    			if (prepend_slot && prepend_slot.p && changed.$$scope) {
    				prepend_slot.p(get_slot_changes(prepend_slot_template, ctx, changed, get_prepend_slot_changes), get_slot_context(prepend_slot_template, ctx, get_prepend_slot_context));
    			}

    			if (!current || changed.pClasses && div1_class_value !== (div1_class_value = "" + (null_to_empty(ctx.pClasses) + " svelte-8az9n8"))) {
    				attr_dev(div1, "class", div1_class_value);
    			}

    			if (ctx.prepend) {
    				if (if_block2) {
    					if_block2.p(changed, ctx);
    					transition_in(if_block2, 1);
    				} else {
    					if_block2 = create_if_block_1$2(ctx);
    					if_block2.c();
    					transition_in(if_block2, 1);
    					if_block2.m(div4, t6);
    				}
    			} else if (if_block2) {
    				group_outros();

    				transition_out(if_block2, 1, 1, () => {
    					if_block2 = null;
    				});

    				check_outros();
    			}

    			if (!current || changed.focused && div2_class_value !== (div2_class_value = "mx-auto w-0 " + (ctx.focused ? ctx.bg() : "") + " svelte-8az9n8")) {
    				attr_dev(div2, "class", div2_class_value);
    			}

    			if (changed.focused || changed.focused || changed.error) {
    				toggle_class(div2, "w-full", ctx.focused || ctx.error);
    			}

    			if (changed.focused || changed.error) {
    				toggle_class(div2, "bg-error-500", ctx.error);
    			}

    			if (changed.noUnderline || changed.outlined) {
    				toggle_class(div3, "hidden", ctx.noUnderline || ctx.outlined);
    			}

    			if (changed.error) {
    				toggle_class(div4, "text-error-500", ctx.error);
    			}

    			if (ctx.showHint) {
    				if (if_block3) {
    					if_block3.p(changed, ctx);
    					transition_in(if_block3, 1);
    				} else {
    					if_block3 = create_if_block$4(ctx);
    					if_block3.c();
    					transition_in(if_block3, 1);
    					if_block3.m(div5, null);
    				}
    			} else if (if_block3) {
    				group_outros();

    				transition_out(if_block3, 1, 1, () => {
    					if_block3 = null;
    				});

    				check_outros();
    			}

    			if (!current || changed.wClasses && div5_class_value !== (div5_class_value = "" + (null_to_empty(ctx.wClasses) + " svelte-8az9n8"))) {
    				attr_dev(div5, "class", div5_class_value);
    			}
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(append_slot, local);
    			transition_in(if_block0);
    			transition_in(prepend_slot, local);
    			transition_in(if_block2);
    			transition_in(if_block3);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(append_slot, local);
    			transition_out(if_block0);
    			transition_out(prepend_slot, local);
    			transition_out(if_block2);
    			transition_out(if_block3);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div5);
    			if (append_slot) append_slot.d(detaching);
    			if (if_block0) if_block0.d();

    			if (if_block1) {
    				if_block1.d();
    			}

    			if (prepend_slot) prepend_slot.d(detaching);
    			if (if_block2) if_block2.d();
    			if (if_block3) if_block3.d();
    			dispose();
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$8.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    let appendDefault = "absolute right-0 top-0 pb-2 pr-4 pt-4 pointer-events-none text-gray-700";
    let prependDefault = "absolute left-0 top-0 pointer-events-none text-xs text-gray-700";

    function instance$5($$self, $$props, $$invalidate) {
    	let { class: className = "" } = $$props;
    	let { outlined = false } = $$props;
    	let { value = null } = $$props;
    	let { label = "" } = $$props;
    	let { placeholder = "" } = $$props;
    	let { hint = "" } = $$props;
    	let { error = false } = $$props;
    	let { append = "" } = $$props;
    	let { prepend = "" } = $$props;
    	let { persistentHint = false } = $$props;
    	let { textarea = false } = $$props;
    	let { rows = 5 } = $$props;
    	let { select = false } = $$props;
    	let { autocomplete = false } = $$props;
    	let { noUnderline = false } = $$props;
    	let { appendReverse = false } = $$props;
    	let { prependReverse = false } = $$props;
    	let { color = "primary" } = $$props;
    	let { bgColor = "white" } = $$props;
    	let { iconClass = "" } = $$props;
    	let labelDefault = `pt-4 absolute top-0 label-transition block pb-2 px-4 pointer-events-none cursor-text`;
    	let inputDefault = `transition pb-2 pt-6 px-4 rounded-t text-black w-full`;
    	let wrapperDefault = "mt-2 relative pb-6 text-gray-600" + (select || autocomplete ? " select" : "");
    	let { add = "" } = $$props;
    	let { remove = "" } = $$props;
    	let { replace = "" } = $$props;
    	let { labelClasses = labelDefault } = $$props;
    	let { inputClasses = inputDefault } = $$props;
    	let { wrapperClasses = wrapperDefault } = $$props;
    	let { appendClasses = appendDefault } = $$props;
    	let { prependClasses = prependDefault } = $$props;
    	const { bg, border, txt, caret } = utils(color);
    	const l = new ClassBuilder(labelClasses, labelDefault);
    	const i = new ClassBuilder(inputClasses, inputDefault);
    	const w = new ClassBuilder(wrapperClasses, wrapperDefault);
    	const a = new ClassBuilder(appendClasses, appendDefault);
    	const p = new ClassBuilder(prependClasses, prependDefault);
    	let focused = false;
    	let lClasses = "";
    	let iClasses = "";
    	let wClasses = "";
    	let aClasses = "";
    	let pClasses = "";

    	function toggleFocused() {
    		$$invalidate("focused", focused = !focused);
    	}

    	const props = filterProps(
    		[
    			"outlined",
    			"label",
    			"placeholder",
    			"hint",
    			"error",
    			"append",
    			"prepend",
    			"persistentHint",
    			"textarea",
    			"rows",
    			"select",
    			"autocomplete",
    			"noUnderline",
    			"appendReverse",
    			"prependReverse",
    			"color",
    			"bgColor"
    		],
    		$$props
    	);

    	let { $$slots = {}, $$scope } = $$props;

    	function blur_handler(event) {
    		bubble($$self, event);
    	}

    	function change_handler(event) {
    		bubble($$self, event);
    	}

    	function input_handler(event) {
    		bubble($$self, event);
    	}

    	function click_handler(event) {
    		bubble($$self, event);
    	}

    	function focus_handler(event) {
    		bubble($$self, event);
    	}

    	function change_handler_1(event) {
    		bubble($$self, event);
    	}

    	function input_handler_1(event) {
    		bubble($$self, event);
    	}

    	function click_handler_1(event) {
    		bubble($$self, event);
    	}

    	function focus_handler_1(event) {
    		bubble($$self, event);
    	}

    	function blur_handler_1(event) {
    		bubble($$self, event);
    	}

    	function change_handler_2(event) {
    		bubble($$self, event);
    	}

    	function input_handler_2(event) {
    		bubble($$self, event);
    	}

    	function click_handler_2(event) {
    		bubble($$self, event);
    	}

    	function blur_handler_2(event) {
    		bubble($$self, event);
    	}

    	function focus_handler_2(event) {
    		bubble($$self, event);
    	}

    	const click_handler_3 = () => select ? $$invalidate("focused", focused = false) : null;

    	function input_input_handler() {
    		value = this.value;
    		$$invalidate("value", value);
    	}

    	function textarea_1_input_handler() {
    		value = this.value;
    		$$invalidate("value", value);
    	}

    	$$self.$set = $$new_props => {
    		$$invalidate("$$props", $$props = assign(assign({}, $$props), $$new_props));
    		if ("class" in $$new_props) $$invalidate("className", className = $$new_props.class);
    		if ("outlined" in $$new_props) $$invalidate("outlined", outlined = $$new_props.outlined);
    		if ("value" in $$new_props) $$invalidate("value", value = $$new_props.value);
    		if ("label" in $$new_props) $$invalidate("label", label = $$new_props.label);
    		if ("placeholder" in $$new_props) $$invalidate("placeholder", placeholder = $$new_props.placeholder);
    		if ("hint" in $$new_props) $$invalidate("hint", hint = $$new_props.hint);
    		if ("error" in $$new_props) $$invalidate("error", error = $$new_props.error);
    		if ("append" in $$new_props) $$invalidate("append", append = $$new_props.append);
    		if ("prepend" in $$new_props) $$invalidate("prepend", prepend = $$new_props.prepend);
    		if ("persistentHint" in $$new_props) $$invalidate("persistentHint", persistentHint = $$new_props.persistentHint);
    		if ("textarea" in $$new_props) $$invalidate("textarea", textarea = $$new_props.textarea);
    		if ("rows" in $$new_props) $$invalidate("rows", rows = $$new_props.rows);
    		if ("select" in $$new_props) $$invalidate("select", select = $$new_props.select);
    		if ("autocomplete" in $$new_props) $$invalidate("autocomplete", autocomplete = $$new_props.autocomplete);
    		if ("noUnderline" in $$new_props) $$invalidate("noUnderline", noUnderline = $$new_props.noUnderline);
    		if ("appendReverse" in $$new_props) $$invalidate("appendReverse", appendReverse = $$new_props.appendReverse);
    		if ("prependReverse" in $$new_props) $$invalidate("prependReverse", prependReverse = $$new_props.prependReverse);
    		if ("color" in $$new_props) $$invalidate("color", color = $$new_props.color);
    		if ("bgColor" in $$new_props) $$invalidate("bgColor", bgColor = $$new_props.bgColor);
    		if ("iconClass" in $$new_props) $$invalidate("iconClass", iconClass = $$new_props.iconClass);
    		if ("add" in $$new_props) $$invalidate("add", add = $$new_props.add);
    		if ("remove" in $$new_props) $$invalidate("remove", remove = $$new_props.remove);
    		if ("replace" in $$new_props) $$invalidate("replace", replace = $$new_props.replace);
    		if ("labelClasses" in $$new_props) $$invalidate("labelClasses", labelClasses = $$new_props.labelClasses);
    		if ("inputClasses" in $$new_props) $$invalidate("inputClasses", inputClasses = $$new_props.inputClasses);
    		if ("wrapperClasses" in $$new_props) $$invalidate("wrapperClasses", wrapperClasses = $$new_props.wrapperClasses);
    		if ("appendClasses" in $$new_props) $$invalidate("appendClasses", appendClasses = $$new_props.appendClasses);
    		if ("prependClasses" in $$new_props) $$invalidate("prependClasses", prependClasses = $$new_props.prependClasses);
    		if ("$$scope" in $$new_props) $$invalidate("$$scope", $$scope = $$new_props.$$scope);
    	};

    	$$self.$capture_state = () => {
    		return {
    			className,
    			outlined,
    			value,
    			label,
    			placeholder,
    			hint,
    			error,
    			append,
    			prepend,
    			persistentHint,
    			textarea,
    			rows,
    			select,
    			autocomplete,
    			noUnderline,
    			appendReverse,
    			prependReverse,
    			color,
    			bgColor,
    			iconClass,
    			labelDefault,
    			inputDefault,
    			wrapperDefault,
    			appendDefault,
    			prependDefault,
    			add,
    			remove,
    			replace,
    			labelClasses,
    			inputClasses,
    			wrapperClasses,
    			appendClasses,
    			prependClasses,
    			focused,
    			lClasses,
    			iClasses,
    			wClasses,
    			aClasses,
    			pClasses,
    			showHint,
    			labelOnTop
    		};
    	};

    	$$self.$inject_state = $$new_props => {
    		$$invalidate("$$props", $$props = assign(assign({}, $$props), $$new_props));
    		if ("className" in $$props) $$invalidate("className", className = $$new_props.className);
    		if ("outlined" in $$props) $$invalidate("outlined", outlined = $$new_props.outlined);
    		if ("value" in $$props) $$invalidate("value", value = $$new_props.value);
    		if ("label" in $$props) $$invalidate("label", label = $$new_props.label);
    		if ("placeholder" in $$props) $$invalidate("placeholder", placeholder = $$new_props.placeholder);
    		if ("hint" in $$props) $$invalidate("hint", hint = $$new_props.hint);
    		if ("error" in $$props) $$invalidate("error", error = $$new_props.error);
    		if ("append" in $$props) $$invalidate("append", append = $$new_props.append);
    		if ("prepend" in $$props) $$invalidate("prepend", prepend = $$new_props.prepend);
    		if ("persistentHint" in $$props) $$invalidate("persistentHint", persistentHint = $$new_props.persistentHint);
    		if ("textarea" in $$props) $$invalidate("textarea", textarea = $$new_props.textarea);
    		if ("rows" in $$props) $$invalidate("rows", rows = $$new_props.rows);
    		if ("select" in $$props) $$invalidate("select", select = $$new_props.select);
    		if ("autocomplete" in $$props) $$invalidate("autocomplete", autocomplete = $$new_props.autocomplete);
    		if ("noUnderline" in $$props) $$invalidate("noUnderline", noUnderline = $$new_props.noUnderline);
    		if ("appendReverse" in $$props) $$invalidate("appendReverse", appendReverse = $$new_props.appendReverse);
    		if ("prependReverse" in $$props) $$invalidate("prependReverse", prependReverse = $$new_props.prependReverse);
    		if ("color" in $$props) $$invalidate("color", color = $$new_props.color);
    		if ("bgColor" in $$props) $$invalidate("bgColor", bgColor = $$new_props.bgColor);
    		if ("iconClass" in $$props) $$invalidate("iconClass", iconClass = $$new_props.iconClass);
    		if ("labelDefault" in $$props) labelDefault = $$new_props.labelDefault;
    		if ("inputDefault" in $$props) inputDefault = $$new_props.inputDefault;
    		if ("wrapperDefault" in $$props) $$invalidate("wrapperDefault", wrapperDefault = $$new_props.wrapperDefault);
    		if ("appendDefault" in $$props) $$invalidate("appendDefault", appendDefault = $$new_props.appendDefault);
    		if ("prependDefault" in $$props) $$invalidate("prependDefault", prependDefault = $$new_props.prependDefault);
    		if ("add" in $$props) $$invalidate("add", add = $$new_props.add);
    		if ("remove" in $$props) $$invalidate("remove", remove = $$new_props.remove);
    		if ("replace" in $$props) $$invalidate("replace", replace = $$new_props.replace);
    		if ("labelClasses" in $$props) $$invalidate("labelClasses", labelClasses = $$new_props.labelClasses);
    		if ("inputClasses" in $$props) $$invalidate("inputClasses", inputClasses = $$new_props.inputClasses);
    		if ("wrapperClasses" in $$props) $$invalidate("wrapperClasses", wrapperClasses = $$new_props.wrapperClasses);
    		if ("appendClasses" in $$props) $$invalidate("appendClasses", appendClasses = $$new_props.appendClasses);
    		if ("prependClasses" in $$props) $$invalidate("prependClasses", prependClasses = $$new_props.prependClasses);
    		if ("focused" in $$props) $$invalidate("focused", focused = $$new_props.focused);
    		if ("lClasses" in $$props) $$invalidate("lClasses", lClasses = $$new_props.lClasses);
    		if ("iClasses" in $$props) $$invalidate("iClasses", iClasses = $$new_props.iClasses);
    		if ("wClasses" in $$props) $$invalidate("wClasses", wClasses = $$new_props.wClasses);
    		if ("aClasses" in $$props) $$invalidate("aClasses", aClasses = $$new_props.aClasses);
    		if ("pClasses" in $$props) $$invalidate("pClasses", pClasses = $$new_props.pClasses);
    		if ("showHint" in $$props) $$invalidate("showHint", showHint = $$new_props.showHint);
    		if ("labelOnTop" in $$props) $$invalidate("labelOnTop", labelOnTop = $$new_props.labelOnTop);
    	};

    	let showHint;
    	let labelOnTop;

    	$$self.$$.update = (changed = { error: 1, persistentHint: 1, hint: 1, focused: 1, placeholder: 1, value: 1, labelOnTop: 1, outlined: 1, bgColor: 1, prepend: 1, className: 1, add: 1, remove: 1, replace: 1, wrapperClasses: 1, wrapperDefault: 1, appendClasses: 1, appendDefault: 1, prependClasses: 1, prependDefault: 1 }) => {
    		if (changed.error || changed.persistentHint || changed.hint || changed.focused) {
    			 $$invalidate("showHint", showHint = error || (persistentHint ? hint : focused && hint));
    		}

    		if (changed.placeholder || changed.focused || changed.value) {
    			 $$invalidate("labelOnTop", labelOnTop = placeholder || focused || value);
    		}

    		if (changed.focused || changed.error || changed.labelOnTop || changed.outlined || changed.bgColor || changed.prepend) {
    			 $$invalidate("lClasses", lClasses = l.flush().add(txt(), focused && !error).add("label-top text-xs", labelOnTop).remove("pt-4 pb-2 px-4 px-1 pt-0", labelOnTop && outlined).add(`ml-3 p-1 pt-0 mt-0 bg-${bgColor}`, labelOnTop && outlined).remove("px-4", prepend).add("pr-4 pl-6", prepend).get());
    		}

    		if (changed.className || changed.outlined || changed.error || changed.focused || changed.prepend || changed.add || changed.remove || changed.replace) {
    			 $$invalidate("iClasses", iClasses = i.flush().add(className).remove("pt-6 pb-2", outlined).add("border rounded bg-transparent py-4 transition", outlined).add("border-error-500 caret-error-500", error).remove(caret(), error).add(caret(), !error).add(border(), focused && !error).add("border-gray-600", !error && !focused).add("bg-gray-100", !outlined).add("bg-gray-300", focused && !outlined).remove("px-4", prepend).add("pr-4 pl-6", prepend).add(add).remove(remove).replace(replace).get());
    		}

    		if (changed.wrapperClasses || changed.wrapperDefault) {
    			 ($$invalidate("wClasses", wClasses = new ClassBuilder(wrapperClasses, wrapperDefault).get()));
    		}

    		if (changed.appendClasses || changed.appendDefault) {
    			 ($$invalidate("aClasses", aClasses = new ClassBuilder(appendClasses, appendDefault).get()));
    		}

    		if (changed.prependClasses || changed.prependDefault) {
    			 ($$invalidate("pClasses", pClasses = new ClassBuilder(prependClasses, prependDefault).get()));
    		}
    	};

    	return {
    		className,
    		outlined,
    		value,
    		label,
    		placeholder,
    		hint,
    		error,
    		append,
    		prepend,
    		persistentHint,
    		textarea,
    		rows,
    		select,
    		autocomplete,
    		noUnderline,
    		appendReverse,
    		prependReverse,
    		color,
    		bgColor,
    		iconClass,
    		add,
    		remove,
    		replace,
    		labelClasses,
    		inputClasses,
    		wrapperClasses,
    		appendClasses,
    		prependClasses,
    		bg,
    		txt,
    		focused,
    		lClasses,
    		iClasses,
    		wClasses,
    		aClasses,
    		pClasses,
    		toggleFocused,
    		props,
    		showHint,
    		blur_handler,
    		change_handler,
    		input_handler,
    		click_handler,
    		focus_handler,
    		change_handler_1,
    		input_handler_1,
    		click_handler_1,
    		focus_handler_1,
    		blur_handler_1,
    		change_handler_2,
    		input_handler_2,
    		click_handler_2,
    		blur_handler_2,
    		focus_handler_2,
    		click_handler_3,
    		input_input_handler,
    		textarea_1_input_handler,
    		$$props: $$props = exclude_internal_props($$props),
    		$$slots,
    		$$scope
    	};
    }

    class TextField extends SvelteComponentDev {
    	constructor(options) {
    		super(options);

    		init(this, options, instance$5, create_fragment$8, safe_not_equal, {
    			class: "className",
    			outlined: 0,
    			value: 0,
    			label: 0,
    			placeholder: 0,
    			hint: 0,
    			error: 0,
    			append: 0,
    			prepend: 0,
    			persistentHint: 0,
    			textarea: 0,
    			rows: 0,
    			select: 0,
    			autocomplete: 0,
    			noUnderline: 0,
    			appendReverse: 0,
    			prependReverse: 0,
    			color: 0,
    			bgColor: 0,
    			iconClass: 0,
    			add: 0,
    			remove: 0,
    			replace: 0,
    			labelClasses: 0,
    			inputClasses: 0,
    			wrapperClasses: 0,
    			appendClasses: 0,
    			prependClasses: 0
    		});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "TextField",
    			options,
    			id: create_fragment$8.name
    		});
    	}

    	get class() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set class(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get outlined() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set outlined(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get value() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set value(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get label() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set label(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get placeholder() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set placeholder(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get hint() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set hint(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get error() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set error(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get append() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set append(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get prepend() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set prepend(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get persistentHint() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set persistentHint(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get textarea() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set textarea(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get rows() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set rows(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get select() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set select(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get autocomplete() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set autocomplete(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get noUnderline() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set noUnderline(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get appendReverse() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set appendReverse(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get prependReverse() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set prependReverse(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get color() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set color(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get bgColor() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set bgColor(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get iconClass() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set iconClass(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get add() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set add(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get remove() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set remove(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get replace() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set replace(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get labelClasses() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set labelClasses(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get inputClasses() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set inputClasses(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get wrapperClasses() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set wrapperClasses(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get appendClasses() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set appendClasses(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get prependClasses() {
    		throw new Error("<TextField>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set prependClasses(value) {
    		throw new Error("<TextField>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}
    }

    /* node_modules/smelte/src/components/Select/Select.svelte generated by Svelte v3.14.0 */
    const file$9 = "node_modules/smelte/src/components/Select/Select.svelte";
    const get_options_slot_changes = () => ({});
    const get_options_slot_context = () => ({});
    const get_select_slot_changes = () => ({});
    const get_select_slot_context = () => ({});

    // (107:2) {#if showList}
    function create_if_block$5(ctx) {
    	let div;
    	let updating_value;
    	let dispose_options_slot;
    	let current;
    	const options_slot_template = ctx.$$slots.options;
    	const options_slot = create_slot(options_slot_template, ctx, get_options_slot_context);

    	function list_value_binding(value_1) {
    		ctx.list_value_binding.call(null, value_1);
    	}

    	let list_props = {
    		select: true,
    		dense: ctx.dense,
    		items: ctx.filteredItems
    	};

    	if (ctx.value !== void 0) {
    		list_props.value = ctx.value;
    	}

    	const list = new List({ props: list_props, $$inline: true });
    	binding_callbacks.push(() => bind(list, "value", list_value_binding));
    	list.$on("change", ctx.change_handler);

    	const block = {
    		c: function create() {
    			if (!options_slot) {
    				div = element("div");
    				create_component(list.$$.fragment);
    			}

    			if (options_slot) options_slot.c();

    			if (!options_slot) {
    				attr_dev(div, "class", "list");
    				toggle_class(div, "rounded-t-none", !ctx.outlined);
    				add_location(div, file$9, 108, 6, 2632);
    				dispose_options_slot = listen_dev(div, "click", ctx.click_handler_3, false, false, false);
    			}
    		},
    		m: function mount(target, anchor) {
    			if (!options_slot) {
    				insert_dev(target, div, anchor);
    				mount_component(list, div, null);
    			}

    			if (options_slot) {
    				options_slot.m(target, anchor);
    			}

    			current = true;
    		},
    		p: function update(changed, ctx) {
    			if (!options_slot) {
    				const list_changes = {};
    				if (changed.dense) list_changes.dense = ctx.dense;
    				if (changed.filteredItems) list_changes.items = ctx.filteredItems;

    				if (!updating_value && changed.value) {
    					updating_value = true;
    					list_changes.value = ctx.value;
    					add_flush_callback(() => updating_value = false);
    				}

    				list.$set(list_changes);

    				if (changed.outlined) {
    					toggle_class(div, "rounded-t-none", !ctx.outlined);
    				}
    			}

    			if (options_slot && options_slot.p && changed.$$scope) {
    				options_slot.p(get_slot_changes(options_slot_template, ctx, changed, get_options_slot_changes), get_slot_context(options_slot_template, ctx, get_options_slot_context));
    			}
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(list.$$.fragment, local);
    			transition_in(options_slot, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(list.$$.fragment, local);
    			transition_out(options_slot, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (!options_slot) {
    				if (detaching) detach_dev(div);
    				destroy_component(list);
    				dispose_options_slot();
    			}

    			if (options_slot) options_slot.d(detaching);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_if_block$5.name,
    		type: "if",
    		source: "(107:2) {#if showList}",
    		ctx
    	});

    	return block;
    }

    function create_fragment$9(ctx) {
    	let div;
    	let t;
    	let div_class_value;
    	let current;
    	let dispose;
    	const select_slot_template = ctx.$$slots.select;
    	const select_slot = create_slot(select_slot_template, ctx, get_select_slot_context);

    	const textfield_spread_levels = [
    		{ select: true },
    		{ autocomplete: ctx.autocomplete },
    		{ value: ctx.selectedLabel },
    		ctx.props,
    		{ wrapperClasses: ctx.inputWrapperClasses },
    		{ appendClasses: ctx.appendClasses },
    		{ labelClasses: ctx.labelClasses },
    		{ inputClasses: ctx.inputClasses },
    		{ prependClasses: ctx.prependClasses },
    		{ append: "arrow_drop_down" },
    		{ appendReverse: ctx.showList }
    	];

    	let textfield_props = {};

    	for (let i = 0; i < textfield_spread_levels.length; i += 1) {
    		textfield_props = assign(textfield_props, textfield_spread_levels[i]);
    	}

    	const textfield = new TextField({ props: textfield_props, $$inline: true });
    	textfield.$on("click", ctx.click_handler_2);
    	textfield.$on("click", ctx.click_handler);
    	textfield.$on("input", ctx.filterItems);
    	let if_block = ctx.showList && create_if_block$5(ctx);

    	const block = {
    		c: function create() {
    			div = element("div");

    			if (!select_slot) {
    				create_component(textfield.$$.fragment);
    			}

    			if (select_slot) select_slot.c();
    			t = space();
    			if (if_block) if_block.c();
    			attr_dev(div, "class", div_class_value = "" + (ctx.wrapperClasses + " " + ctx.className));
    			add_location(div, file$9, 83, 0, 2083);
    			dispose = listen_dev(window, "click", ctx.click_handler_1, false, false, false);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div, anchor);

    			if (!select_slot) {
    				mount_component(textfield, div, null);
    			}

    			if (select_slot) {
    				select_slot.m(div, null);
    			}

    			append_dev(div, t);
    			if (if_block) if_block.m(div, null);
    			current = true;
    		},
    		p: function update(changed, ctx) {
    			if (!select_slot) {
    				const textfield_changes = changed.autocomplete || changed.selectedLabel || changed.props || changed.inputWrapperClasses || changed.appendClasses || changed.labelClasses || changed.inputClasses || changed.prependClasses || changed.showList
    				? get_spread_update(textfield_spread_levels, [
    						textfield_spread_levels[0],
    						changed.autocomplete && ({ autocomplete: ctx.autocomplete }),
    						changed.selectedLabel && ({ value: ctx.selectedLabel }),
    						changed.props && get_spread_object(ctx.props),
    						changed.inputWrapperClasses && ({ wrapperClasses: ctx.inputWrapperClasses }),
    						changed.appendClasses && ({ appendClasses: ctx.appendClasses }),
    						changed.labelClasses && ({ labelClasses: ctx.labelClasses }),
    						changed.inputClasses && ({ inputClasses: ctx.inputClasses }),
    						changed.prependClasses && ({ prependClasses: ctx.prependClasses }),
    						textfield_spread_levels[9],
    						changed.showList && ({ appendReverse: ctx.showList })
    					])
    				: {};

    				textfield.$set(textfield_changes);
    			}

    			if (select_slot && select_slot.p && changed.$$scope) {
    				select_slot.p(get_slot_changes(select_slot_template, ctx, changed, get_select_slot_changes), get_slot_context(select_slot_template, ctx, get_select_slot_context));
    			}

    			if (ctx.showList) {
    				if (if_block) {
    					if_block.p(changed, ctx);
    					transition_in(if_block, 1);
    				} else {
    					if_block = create_if_block$5(ctx);
    					if_block.c();
    					transition_in(if_block, 1);
    					if_block.m(div, null);
    				}
    			} else if (if_block) {
    				group_outros();

    				transition_out(if_block, 1, 1, () => {
    					if_block = null;
    				});

    				check_outros();
    			}

    			if (!current || (changed.wrapperClasses || changed.className) && div_class_value !== (div_class_value = "" + (ctx.wrapperClasses + " " + ctx.className))) {
    				attr_dev(div, "class", div_class_value);
    			}
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(textfield.$$.fragment, local);
    			transition_in(select_slot, local);
    			transition_in(if_block);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(textfield.$$.fragment, local);
    			transition_out(select_slot, local);
    			transition_out(if_block);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div);

    			if (!select_slot) {
    				destroy_component(textfield);
    			}

    			if (select_slot) select_slot.d(detaching);
    			if (if_block) if_block.d();
    			dispose();
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$9.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    function process(it) {
    	return it.map(i => typeof i !== "object" ? { value: i, text: i } : i);
    }

    function instance$6($$self, $$props, $$invalidate) {
    	let { items = [] } = $$props;
    	let { class: className = "" } = $$props;
    	let { value = "" } = $$props;
    	let { text = "" } = $$props;
    	let { label = "" } = $$props;
    	let { selectedLabel = "" } = $$props;
    	let { color = "primary" } = $$props;
    	let { outlined = false } = $$props;
    	let { placeholder = "" } = $$props;
    	let { hint = "" } = $$props;
    	let { error = false } = $$props;
    	let { append = "" } = $$props;
    	let { dense = false } = $$props;
    	let { persistentHint = false } = $$props;
    	let { autocomplete = false } = $$props;
    	let { noUnderline = false } = $$props;
    	let { wrapperClasses = "cursor-pointer relative pb-4" } = $$props;
    	let { showList = false } = $$props;
    	let { inputWrapperClasses = i => i } = $$props;
    	let { appendClasses = i => i } = $$props;
    	let { labelClasses = i => i } = $$props;
    	let { inputClasses = i => i } = $$props;
    	let { prependClasses = i => i } = $$props;
    	let { add = "" } = $$props;
    	let { remove = "" } = $$props;
    	let { replace = "" } = $$props;
    	let filteredItems = items;
    	let itemsProcessed = [];

    	const props = {
    		outlined,
    		label,
    		placeholder,
    		hint,
    		error,
    		append,
    		persistentHint,
    		color,
    		add,
    		remove,
    		replace,
    		noUnderline
    	};

    	onMount(() => {
    		$$invalidate("selectedLabel", selectedLabel = getLabel(value));
    	});

    	const dispatch = createEventDispatcher();

    	function getLabel(value) {
    		return value
    		? (itemsProcessed.find(i => i.value === value) || ({ text: "" })).text
    		: "";
    	}

    	function filterItems({ target }) {
    		$$invalidate("filteredItems", filteredItems = itemsProcessed.filter(i => i.text.toLowerCase().includes(target.value.toLowerCase())));
    	}

    	const writable_props = [
    		"items",
    		"class",
    		"value",
    		"text",
    		"label",
    		"selectedLabel",
    		"color",
    		"outlined",
    		"placeholder",
    		"hint",
    		"error",
    		"append",
    		"dense",
    		"persistentHint",
    		"autocomplete",
    		"noUnderline",
    		"wrapperClasses",
    		"showList",
    		"inputWrapperClasses",
    		"appendClasses",
    		"labelClasses",
    		"inputClasses",
    		"prependClasses",
    		"add",
    		"remove",
    		"replace"
    	];

    	Object.keys($$props).forEach(key => {
    		if (!writable_props.includes(key) && !key.startsWith("$$")) console.warn(`<Select> was created with unknown prop '${key}'`);
    	});

    	let { $$slots = {}, $$scope } = $$props;
    	const click_handler_1 = () => $$invalidate("showList", showList = false);

    	const click_handler_2 = e => {
    		e.stopPropagation();
    		$$invalidate("showList", showList = true);
    	};

    	function click_handler(event) {
    		bubble($$self, event);
    	}

    	function list_value_binding(value_1) {
    		value = value_1;
    		$$invalidate("value", value);
    	}

    	const change_handler = ({ detail }) => {
    		$$invalidate("selectedLabel", selectedLabel = getLabel(detail));
    		dispatch("change", detail);
    	};

    	const click_handler_3 = () => $$invalidate("showList", showList = false);

    	$$self.$set = $$props => {
    		if ("items" in $$props) $$invalidate("items", items = $$props.items);
    		if ("class" in $$props) $$invalidate("className", className = $$props.class);
    		if ("value" in $$props) $$invalidate("value", value = $$props.value);
    		if ("text" in $$props) $$invalidate("text", text = $$props.text);
    		if ("label" in $$props) $$invalidate("label", label = $$props.label);
    		if ("selectedLabel" in $$props) $$invalidate("selectedLabel", selectedLabel = $$props.selectedLabel);
    		if ("color" in $$props) $$invalidate("color", color = $$props.color);
    		if ("outlined" in $$props) $$invalidate("outlined", outlined = $$props.outlined);
    		if ("placeholder" in $$props) $$invalidate("placeholder", placeholder = $$props.placeholder);
    		if ("hint" in $$props) $$invalidate("hint", hint = $$props.hint);
    		if ("error" in $$props) $$invalidate("error", error = $$props.error);
    		if ("append" in $$props) $$invalidate("append", append = $$props.append);
    		if ("dense" in $$props) $$invalidate("dense", dense = $$props.dense);
    		if ("persistentHint" in $$props) $$invalidate("persistentHint", persistentHint = $$props.persistentHint);
    		if ("autocomplete" in $$props) $$invalidate("autocomplete", autocomplete = $$props.autocomplete);
    		if ("noUnderline" in $$props) $$invalidate("noUnderline", noUnderline = $$props.noUnderline);
    		if ("wrapperClasses" in $$props) $$invalidate("wrapperClasses", wrapperClasses = $$props.wrapperClasses);
    		if ("showList" in $$props) $$invalidate("showList", showList = $$props.showList);
    		if ("inputWrapperClasses" in $$props) $$invalidate("inputWrapperClasses", inputWrapperClasses = $$props.inputWrapperClasses);
    		if ("appendClasses" in $$props) $$invalidate("appendClasses", appendClasses = $$props.appendClasses);
    		if ("labelClasses" in $$props) $$invalidate("labelClasses", labelClasses = $$props.labelClasses);
    		if ("inputClasses" in $$props) $$invalidate("inputClasses", inputClasses = $$props.inputClasses);
    		if ("prependClasses" in $$props) $$invalidate("prependClasses", prependClasses = $$props.prependClasses);
    		if ("add" in $$props) $$invalidate("add", add = $$props.add);
    		if ("remove" in $$props) $$invalidate("remove", remove = $$props.remove);
    		if ("replace" in $$props) $$invalidate("replace", replace = $$props.replace);
    		if ("$$scope" in $$props) $$invalidate("$$scope", $$scope = $$props.$$scope);
    	};

    	$$self.$capture_state = () => {
    		return {
    			items,
    			className,
    			value,
    			text,
    			label,
    			selectedLabel,
    			color,
    			outlined,
    			placeholder,
    			hint,
    			error,
    			append,
    			dense,
    			persistentHint,
    			autocomplete,
    			noUnderline,
    			wrapperClasses,
    			showList,
    			inputWrapperClasses,
    			appendClasses,
    			labelClasses,
    			inputClasses,
    			prependClasses,
    			add,
    			remove,
    			replace,
    			filteredItems,
    			itemsProcessed
    		};
    	};

    	$$self.$inject_state = $$props => {
    		if ("items" in $$props) $$invalidate("items", items = $$props.items);
    		if ("className" in $$props) $$invalidate("className", className = $$props.className);
    		if ("value" in $$props) $$invalidate("value", value = $$props.value);
    		if ("text" in $$props) $$invalidate("text", text = $$props.text);
    		if ("label" in $$props) $$invalidate("label", label = $$props.label);
    		if ("selectedLabel" in $$props) $$invalidate("selectedLabel", selectedLabel = $$props.selectedLabel);
    		if ("color" in $$props) $$invalidate("color", color = $$props.color);
    		if ("outlined" in $$props) $$invalidate("outlined", outlined = $$props.outlined);
    		if ("placeholder" in $$props) $$invalidate("placeholder", placeholder = $$props.placeholder);
    		if ("hint" in $$props) $$invalidate("hint", hint = $$props.hint);
    		if ("error" in $$props) $$invalidate("error", error = $$props.error);
    		if ("append" in $$props) $$invalidate("append", append = $$props.append);
    		if ("dense" in $$props) $$invalidate("dense", dense = $$props.dense);
    		if ("persistentHint" in $$props) $$invalidate("persistentHint", persistentHint = $$props.persistentHint);
    		if ("autocomplete" in $$props) $$invalidate("autocomplete", autocomplete = $$props.autocomplete);
    		if ("noUnderline" in $$props) $$invalidate("noUnderline", noUnderline = $$props.noUnderline);
    		if ("wrapperClasses" in $$props) $$invalidate("wrapperClasses", wrapperClasses = $$props.wrapperClasses);
    		if ("showList" in $$props) $$invalidate("showList", showList = $$props.showList);
    		if ("inputWrapperClasses" in $$props) $$invalidate("inputWrapperClasses", inputWrapperClasses = $$props.inputWrapperClasses);
    		if ("appendClasses" in $$props) $$invalidate("appendClasses", appendClasses = $$props.appendClasses);
    		if ("labelClasses" in $$props) $$invalidate("labelClasses", labelClasses = $$props.labelClasses);
    		if ("inputClasses" in $$props) $$invalidate("inputClasses", inputClasses = $$props.inputClasses);
    		if ("prependClasses" in $$props) $$invalidate("prependClasses", prependClasses = $$props.prependClasses);
    		if ("add" in $$props) $$invalidate("add", add = $$props.add);
    		if ("remove" in $$props) $$invalidate("remove", remove = $$props.remove);
    		if ("replace" in $$props) $$invalidate("replace", replace = $$props.replace);
    		if ("filteredItems" in $$props) $$invalidate("filteredItems", filteredItems = $$props.filteredItems);
    		if ("itemsProcessed" in $$props) itemsProcessed = $$props.itemsProcessed;
    	};

    	$$self.$$.update = (changed = { items: 1 }) => {
    		if (changed.items) {
    			 itemsProcessed = process(items);
    		}
    	};

    	return {
    		items,
    		className,
    		value,
    		text,
    		label,
    		selectedLabel,
    		color,
    		outlined,
    		placeholder,
    		hint,
    		error,
    		append,
    		dense,
    		persistentHint,
    		autocomplete,
    		noUnderline,
    		wrapperClasses,
    		showList,
    		inputWrapperClasses,
    		appendClasses,
    		labelClasses,
    		inputClasses,
    		prependClasses,
    		add,
    		remove,
    		replace,
    		filteredItems,
    		props,
    		dispatch,
    		getLabel,
    		filterItems,
    		click_handler_1,
    		click_handler_2,
    		click_handler,
    		list_value_binding,
    		change_handler,
    		click_handler_3,
    		$$slots,
    		$$scope
    	};
    }

    class Select extends SvelteComponentDev {
    	constructor(options) {
    		super(options);

    		init(this, options, instance$6, create_fragment$9, safe_not_equal, {
    			items: 0,
    			class: "className",
    			value: 0,
    			text: 0,
    			label: 0,
    			selectedLabel: 0,
    			color: 0,
    			outlined: 0,
    			placeholder: 0,
    			hint: 0,
    			error: 0,
    			append: 0,
    			dense: 0,
    			persistentHint: 0,
    			autocomplete: 0,
    			noUnderline: 0,
    			wrapperClasses: 0,
    			showList: 0,
    			inputWrapperClasses: 0,
    			appendClasses: 0,
    			labelClasses: 0,
    			inputClasses: 0,
    			prependClasses: 0,
    			add: 0,
    			remove: 0,
    			replace: 0
    		});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "Select",
    			options,
    			id: create_fragment$9.name
    		});
    	}

    	get items() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set items(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get class() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set class(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get value() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set value(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get text() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set text(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get label() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set label(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get selectedLabel() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set selectedLabel(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get color() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set color(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get outlined() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set outlined(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get placeholder() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set placeholder(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get hint() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set hint(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get error() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set error(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get append() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set append(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get dense() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set dense(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get persistentHint() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set persistentHint(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get autocomplete() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set autocomplete(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get noUnderline() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set noUnderline(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get wrapperClasses() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set wrapperClasses(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get showList() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set showList(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get inputWrapperClasses() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set inputWrapperClasses(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get appendClasses() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set appendClasses(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get labelClasses() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set labelClasses(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get inputClasses() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set inputClasses(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get prependClasses() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set prependClasses(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get add() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set add(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get remove() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set remove(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get replace() {
    		throw new Error("<Select>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set replace(value) {
    		throw new Error("<Select>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}
    }

    /* node_modules/smelte/src/components/ProgressLinear/ProgressLinear.svelte generated by Svelte v3.14.0 */
    const file$a = "node_modules/smelte/src/components/ProgressLinear/ProgressLinear.svelte";

    function create_fragment$a(ctx) {
    	let div2;
    	let div0;
    	let div0_class_value;
    	let div0_style_value;
    	let t;
    	let div1;
    	let div1_class_value;
    	let div2_class_value;
    	let div2_transition;
    	let current;

    	const block = {
    		c: function create() {
    			div2 = element("div");
    			div0 = element("div");
    			t = space();
    			div1 = element("div");
    			attr_dev(div0, "class", div0_class_value = "bg-" + ctx.color + "-500 h-1 absolute" + " svelte-8m92aa");
    			attr_dev(div0, "style", div0_style_value = ctx.progress ? `width: ${ctx.progress}%` : "");
    			toggle_class(div0, "inc", !ctx.progress);
    			toggle_class(div0, "transition", ctx.progress);
    			add_location(div0, file$a, 87, 2, 2789);
    			attr_dev(div1, "class", div1_class_value = "bg-" + ctx.color + "-500 h-1 absolute dec" + " svelte-8m92aa");
    			toggle_class(div1, "hidden", ctx.progress);
    			add_location(div1, file$a, 92, 2, 2947);
    			attr_dev(div2, "class", div2_class_value = "top-0 left-0 w-full h-1 bg-" + ctx.color + "-100 overflow-hidden relative" + " svelte-8m92aa");
    			toggle_class(div2, "fixed", ctx.app);
    			toggle_class(div2, "z-50", ctx.app);
    			toggle_class(div2, "hidden", ctx.app && !ctx.initialized);
    			add_location(div2, file$a, 81, 0, 2592);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div2, anchor);
    			append_dev(div2, div0);
    			append_dev(div2, t);
    			append_dev(div2, div1);
    			current = true;
    		},
    		p: function update(changed, ctx) {
    			if (!current || changed.color && div0_class_value !== (div0_class_value = "bg-" + ctx.color + "-500 h-1 absolute" + " svelte-8m92aa")) {
    				attr_dev(div0, "class", div0_class_value);
    			}

    			if (!current || changed.progress && div0_style_value !== (div0_style_value = ctx.progress ? `width: ${ctx.progress}%` : "")) {
    				attr_dev(div0, "style", div0_style_value);
    			}

    			if (changed.color || changed.progress) {
    				toggle_class(div0, "inc", !ctx.progress);
    			}

    			if (changed.color || changed.progress) {
    				toggle_class(div0, "transition", ctx.progress);
    			}

    			if (!current || changed.color && div1_class_value !== (div1_class_value = "bg-" + ctx.color + "-500 h-1 absolute dec" + " svelte-8m92aa")) {
    				attr_dev(div1, "class", div1_class_value);
    			}

    			if (changed.color || changed.progress) {
    				toggle_class(div1, "hidden", ctx.progress);
    			}

    			if (!current || changed.color && div2_class_value !== (div2_class_value = "top-0 left-0 w-full h-1 bg-" + ctx.color + "-100 overflow-hidden relative" + " svelte-8m92aa")) {
    				attr_dev(div2, "class", div2_class_value);
    			}

    			if (changed.color || changed.app) {
    				toggle_class(div2, "fixed", ctx.app);
    			}

    			if (changed.color || changed.app) {
    				toggle_class(div2, "z-50", ctx.app);
    			}

    			if (changed.color || changed.app || changed.initialized) {
    				toggle_class(div2, "hidden", ctx.app && !ctx.initialized);
    			}
    		},
    		i: function intro(local) {
    			if (current) return;

    			add_render_callback(() => {
    				if (!div2_transition) div2_transition = create_bidirectional_transition(div2, slide, { duration: 300 }, true);
    				div2_transition.run(1);
    			});

    			current = true;
    		},
    		o: function outro(local) {
    			if (!div2_transition) div2_transition = create_bidirectional_transition(div2, slide, { duration: 300 }, false);
    			div2_transition.run(0);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div2);
    			if (detaching && div2_transition) div2_transition.end();
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$a.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    function instance$7($$self, $$props, $$invalidate) {
    	let { app = false } = $$props;
    	let { progress = 0 } = $$props;
    	let { color = "primary" } = $$props;
    	let initialized = false;

    	onMount(() => {
    		if (!app) return;

    		setTimeout(
    			() => {
    				$$invalidate("initialized", initialized = true);
    			},
    			200
    		);
    	});

    	const writable_props = ["app", "progress", "color"];

    	Object.keys($$props).forEach(key => {
    		if (!writable_props.includes(key) && !key.startsWith("$$")) console.warn(`<ProgressLinear> was created with unknown prop '${key}'`);
    	});

    	$$self.$set = $$props => {
    		if ("app" in $$props) $$invalidate("app", app = $$props.app);
    		if ("progress" in $$props) $$invalidate("progress", progress = $$props.progress);
    		if ("color" in $$props) $$invalidate("color", color = $$props.color);
    	};

    	$$self.$capture_state = () => {
    		return { app, progress, color, initialized };
    	};

    	$$self.$inject_state = $$props => {
    		if ("app" in $$props) $$invalidate("app", app = $$props.app);
    		if ("progress" in $$props) $$invalidate("progress", progress = $$props.progress);
    		if ("color" in $$props) $$invalidate("color", color = $$props.color);
    		if ("initialized" in $$props) $$invalidate("initialized", initialized = $$props.initialized);
    	};

    	return { app, progress, color, initialized };
    }

    class ProgressLinear extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, instance$7, create_fragment$a, safe_not_equal, { app: 0, progress: 0, color: 0 });

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "ProgressLinear",
    			options,
    			id: create_fragment$a.name
    		});
    	}

    	get app() {
    		throw new Error("<ProgressLinear>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set app(value) {
    		throw new Error("<ProgressLinear>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get progress() {
    		throw new Error("<ProgressLinear>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set progress(value) {
    		throw new Error("<ProgressLinear>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get color() {
    		throw new Error("<ProgressLinear>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set color(value) {
    		throw new Error("<ProgressLinear>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}
    }

    function sort(data, col, asc) {
      if (!col) return data;

      if (col.sort) return col.sort(data);

      const sorted = data.sort((a, b) => {
        const valA = col.value ? col.value(a) : a[col.field];
        const valB = col.value ? col.value(b) : b[col.field];

        const first = asc ? valA : valB;
        const second = asc ? valB : valA;

        if (typeof valA === "number") {
          return first - second;
        }

        return ("" + first).localeCompare(second);
      });

      return sorted;
    }

    /* node_modules/smelte/src/components/DataTable/DataTable.svelte generated by Svelte v3.14.0 */

    const { Object: Object_1 } = globals;
    const file$b = "node_modules/smelte/src/components/DataTable/DataTable.svelte";
    const get_footer_slot_changes = () => ({});
    const get_footer_slot_context = () => ({});
    const get_pagination_slot_changes = () => ({});
    const get_pagination_slot_context = () => ({});
    const get_edit_dialog_slot_changes = () => ({});
    const get_edit_dialog_slot_context = () => ({});

    function get_each_context_1$1(ctx, list, i) {
    	const child_ctx = Object_1.create(ctx);
    	child_ctx.column = list[i];
    	child_ctx.i = i;
    	return child_ctx;
    }

    const get_item_slot_changes$1 = () => ({});
    const get_item_slot_context$1 = () => ({});

    function get_each_context$2(ctx, list, i) {
    	const child_ctx = Object_1.create(ctx);
    	child_ctx.item = list[i];
    	child_ctx.j = i;
    	return child_ctx;
    }

    const get_header_slot_changes = () => ({});
    const get_header_slot_context = () => ({});

    function get_each_context_2(ctx, list, i) {
    	const child_ctx = Object_1.create(ctx);
    	child_ctx.column = list[i];
    	child_ctx.i = i;
    	return child_ctx;
    }

    // (138:12) {#if sortable && column.sortable !== false}
    function create_if_block_4$1(ctx) {
    	let span;
    	let current;

    	const icon = new Icon({
    			props: {
    				small: true,
    				color: "text-gray-400",
    				$$slots: { default: [create_default_slot$4] },
    				$$scope: { ctx }
    			},
    			$$inline: true
    		});

    	const block = {
    		c: function create() {
    			span = element("span");
    			create_component(icon.$$.fragment);
    			attr_dev(span, "class", "sort svelte-13j4zi7");
    			toggle_class(span, "asc", !ctx.asc && ctx.sortBy === ctx.column);
    			add_location(span, file$b, 138, 14, 5011);
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, span, anchor);
    			mount_component(icon, span, null);
    			current = true;
    		},
    		p: function update(changed, ctx) {
    			const icon_changes = {};

    			if (changed.$$scope) {
    				icon_changes.$$scope = { changed, ctx };
    			}

    			icon.$set(icon_changes);

    			if (changed.asc || changed.sortBy || changed.columns) {
    				toggle_class(span, "asc", !ctx.asc && ctx.sortBy === ctx.column);
    			}
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(icon.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(icon.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(span);
    			destroy_component(icon);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_if_block_4$1.name,
    		type: "if",
    		source: "(138:12) {#if sortable && column.sortable !== false}",
    		ctx
    	});

    	return block;
    }

    // (140:16) <Icon small color="text-gray-400">
    function create_default_slot$4(ctx) {
    	let t;

    	const block = {
    		c: function create() {
    			t = text("arrow_downward");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, t, anchor);
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(t);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_default_slot$4.name,
    		type: "slot",
    		source: "(140:16) <Icon small color=\\\"text-gray-400\\\">",
    		ctx
    	});

    	return block;
    }

    // (123:4) {#each columns as column, i}
    function create_each_block_2(ctx) {
    	let th;
    	let div;
    	let t0;
    	let span;
    	let t1_value = (ctx.column.label || ctx.column.field) + "";
    	let t1;
    	let t2;
    	let dispose_header_slot;
    	let current;
    	const header_slot_template = ctx.$$slots.header;
    	const header_slot = create_slot(header_slot_template, ctx, get_header_slot_context);
    	let if_block = ctx.sortable && ctx.column.sortable !== false && create_if_block_4$1(ctx);

    	function click_handler(...args) {
    		return ctx.click_handler(ctx, ...args);
    	}

    	const block = {
    		c: function create() {
    			if (!header_slot) {
    				th = element("th");
    				div = element("div");
    				if (if_block) if_block.c();
    				t0 = space();
    				span = element("span");
    				t1 = text(t1_value);
    				t2 = space();
    			}

    			if (header_slot) header_slot.c();

    			if (!header_slot) {
    				add_location(span, file$b, 142, 12, 5193);
    				attr_dev(div, "class", "sort-wrapper svelte-13j4zi7");
    				add_location(div, file$b, 136, 10, 4914);
    				attr_dev(th, "class", "capitalize svelte-13j4zi7");
    				toggle_class(th, "cursor-pointer", ctx.sortable || ctx.column.sortable);
    				add_location(th, file$b, 124, 8, 4546);
    				dispose_header_slot = listen_dev(th, "click", click_handler, false, false, false);
    			}
    		},
    		m: function mount(target, anchor) {
    			if (!header_slot) {
    				insert_dev(target, th, anchor);
    				append_dev(th, div);
    				if (if_block) if_block.m(div, null);
    				append_dev(div, t0);
    				append_dev(div, span);
    				append_dev(span, t1);
    				insert_dev(target, t2, anchor);
    			}

    			if (header_slot) {
    				header_slot.m(target, anchor);
    			}

    			current = true;
    		},
    		p: function update(changed, new_ctx) {
    			ctx = new_ctx;

    			if (!header_slot) {
    				if (ctx.sortable && ctx.column.sortable !== false) {
    					if (if_block) {
    						if_block.p(changed, ctx);
    						transition_in(if_block, 1);
    					} else {
    						if_block = create_if_block_4$1(ctx);
    						if_block.c();
    						transition_in(if_block, 1);
    						if_block.m(div, t0);
    					}
    				} else if (if_block) {
    					group_outros();

    					transition_out(if_block, 1, 1, () => {
    						if_block = null;
    					});

    					check_outros();
    				}

    				if ((!current || changed.columns) && t1_value !== (t1_value = (ctx.column.label || ctx.column.field) + "")) set_data_dev(t1, t1_value);

    				if (changed.sortable || changed.columns) {
    					toggle_class(th, "cursor-pointer", ctx.sortable || ctx.column.sortable);
    				}
    			}

    			if (header_slot && header_slot.p && changed.$$scope) {
    				header_slot.p(get_slot_changes(header_slot_template, ctx, changed, get_header_slot_changes), get_slot_context(header_slot_template, ctx, get_header_slot_context));
    			}
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(if_block);
    			transition_in(header_slot, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(if_block);
    			transition_out(header_slot, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (!header_slot) {
    				if (detaching) detach_dev(th);
    				if (if_block) if_block.d();
    				if (detaching) detach_dev(t2);
    				dispose_header_slot();
    			}

    			if (header_slot) header_slot.d(detaching);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_each_block_2.name,
    		type: "each",
    		source: "(123:4) {#each columns as column, i}",
    		ctx
    	});

    	return block;
    }

    // (149:2) {#if loading && !hideProgress}
    function create_if_block_3$1(ctx) {
    	let div;
    	let div_transition;
    	let current;
    	const progresslinear = new ProgressLinear({ $$inline: true });

    	const block = {
    		c: function create() {
    			div = element("div");
    			create_component(progresslinear.$$.fragment);
    			attr_dev(div, "class", "absolute w-full");
    			add_location(div, file$b, 149, 4, 5342);
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div, anchor);
    			mount_component(progresslinear, div, null);
    			current = true;
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(progresslinear.$$.fragment, local);

    			add_render_callback(() => {
    				if (!div_transition) div_transition = create_bidirectional_transition(div, slide, {}, true);
    				div_transition.run(1);
    			});

    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(progresslinear.$$.fragment, local);
    			if (!div_transition) div_transition = create_bidirectional_transition(div, slide, {}, false);
    			div_transition.run(0);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div);
    			destroy_component(progresslinear);
    			if (detaching && div_transition) div_transition.end();
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_if_block_3$1.name,
    		type: "if",
    		source: "(149:2) {#if loading && !hideProgress}",
    		ctx
    	});

    	return block;
    }

    // (169:14) {#if editable && column.editable !== false && editing[j] === i}
    function create_if_block_2$2(ctx) {
    	let div;
    	let current;
    	const edit_dialog_slot_template = ctx.$$slots["edit-dialog"];
    	const edit_dialog_slot = create_slot(edit_dialog_slot_template, ctx, get_edit_dialog_slot_context);

    	function blur_handler(...args) {
    		return ctx.blur_handler(ctx, ...args);
    	}

    	const textfield = new TextField({
    			props: {
    				value: ctx.item[ctx.column.field],
    				textarea: ctx.column.textarea,
    				remove: "bg-gray-100 bg-gray-300"
    			},
    			$$inline: true
    		});

    	textfield.$on("change", ctx.change_handler);
    	textfield.$on("blur", blur_handler);

    	const block = {
    		c: function create() {
    			if (!edit_dialog_slot) {
    				div = element("div");
    				create_component(textfield.$$.fragment);
    			}

    			if (edit_dialog_slot) edit_dialog_slot.c();

    			if (!edit_dialog_slot) {
    				attr_dev(div, "class", "absolute left-0 top-0 z-10 bg-white p-2 elevation-3 rounded");
    				set_style(div, "width", "300px");
    				add_location(div, file$b, 170, 18, 6050);
    			}
    		},
    		m: function mount(target, anchor) {
    			if (!edit_dialog_slot) {
    				insert_dev(target, div, anchor);
    				mount_component(textfield, div, null);
    			}

    			if (edit_dialog_slot) {
    				edit_dialog_slot.m(target, anchor);
    			}

    			current = true;
    		},
    		p: function update(changed, new_ctx) {
    			ctx = new_ctx;

    			if (!edit_dialog_slot) {
    				const textfield_changes = {};
    				if (changed.sorted || changed.columns) textfield_changes.value = ctx.item[ctx.column.field];
    				if (changed.columns) textfield_changes.textarea = ctx.column.textarea;
    				textfield.$set(textfield_changes);
    			}

    			if (edit_dialog_slot && edit_dialog_slot.p && changed.$$scope) {
    				edit_dialog_slot.p(get_slot_changes(edit_dialog_slot_template, ctx, changed, get_edit_dialog_slot_changes), get_slot_context(edit_dialog_slot_template, ctx, get_edit_dialog_slot_context));
    			}
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(textfield.$$.fragment, local);
    			transition_in(edit_dialog_slot, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(textfield.$$.fragment, local);
    			transition_out(edit_dialog_slot, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (!edit_dialog_slot) {
    				if (detaching) detach_dev(div);
    				destroy_component(textfield);
    			}

    			if (edit_dialog_slot) edit_dialog_slot.d(detaching);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_if_block_2$2.name,
    		type: "if",
    		source: "(169:14) {#if editable && column.editable !== false && editing[j] === i}",
    		ctx
    	});

    	return block;
    }

    // (191:14) {:else}
    function create_else_block$2(ctx) {
    	let html_tag;
    	let raw_value = ctx.item[ctx.column.field] + "";

    	const block = {
    		c: function create() {
    			html_tag = new HtmlTag(raw_value, null);
    		},
    		m: function mount(target, anchor) {
    			html_tag.m(target, anchor);
    		},
    		p: function update(changed, ctx) {
    			if ((changed.sorted || changed.columns) && raw_value !== (raw_value = ctx.item[ctx.column.field] + "")) html_tag.p(raw_value);
    		},
    		d: function destroy(detaching) {
    			if (detaching) html_tag.d();
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_else_block$2.name,
    		type: "else",
    		source: "(191:14) {:else}",
    		ctx
    	});

    	return block;
    }

    // (189:14) {#if column.value}
    function create_if_block_1$3(ctx) {
    	let html_tag;
    	let raw_value = ctx.column.value(ctx.item) + "";

    	const block = {
    		c: function create() {
    			html_tag = new HtmlTag(raw_value, null);
    		},
    		m: function mount(target, anchor) {
    			html_tag.m(target, anchor);
    		},
    		p: function update(changed, ctx) {
    			if ((changed.columns || changed.sorted) && raw_value !== (raw_value = ctx.column.value(ctx.item) + "")) html_tag.p(raw_value);
    		},
    		d: function destroy(detaching) {
    			if (detaching) html_tag.d();
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_if_block_1$3.name,
    		type: "if",
    		source: "(189:14) {#if column.value}",
    		ctx
    	});

    	return block;
    }

    // (164:10) {#each columns as column, i}
    function create_each_block_1$1(ctx) {
    	let td;
    	let t0;
    	let t1;
    	let td_class_value;
    	let current;
    	let if_block0 = ctx.editable && ctx.column.editable !== false && ctx.editing[ctx.j] === ctx.i && create_if_block_2$2(ctx);

    	function select_block_type(changed, ctx) {
    		if (ctx.column.value) return create_if_block_1$3;
    		return create_else_block$2;
    	}

    	let current_block_type = select_block_type(null, ctx);
    	let if_block1 = current_block_type(ctx);

    	const block = {
    		c: function create() {
    			td = element("td");
    			if (if_block0) if_block0.c();
    			t0 = space();
    			if_block1.c();
    			t1 = space();
    			attr_dev(td, "class", td_class_value = "relative " + ctx.column.class + " svelte-13j4zi7");
    			toggle_class(td, "cursor-pointer", ctx.editable && ctx.column.editable !== false);
    			add_location(td, file$b, 164, 12, 5773);
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, td, anchor);
    			if (if_block0) if_block0.m(td, null);
    			append_dev(td, t0);
    			if_block1.m(td, null);
    			append_dev(td, t1);
    			current = true;
    		},
    		p: function update(changed, ctx) {
    			if (ctx.editable && ctx.column.editable !== false && ctx.editing[ctx.j] === ctx.i) {
    				if (if_block0) {
    					if_block0.p(changed, ctx);
    					transition_in(if_block0, 1);
    				} else {
    					if_block0 = create_if_block_2$2(ctx);
    					if_block0.c();
    					transition_in(if_block0, 1);
    					if_block0.m(td, t0);
    				}
    			} else if (if_block0) {
    				group_outros();

    				transition_out(if_block0, 1, 1, () => {
    					if_block0 = null;
    				});

    				check_outros();
    			}

    			if (current_block_type === (current_block_type = select_block_type(changed, ctx)) && if_block1) {
    				if_block1.p(changed, ctx);
    			} else {
    				if_block1.d(1);
    				if_block1 = current_block_type(ctx);

    				if (if_block1) {
    					if_block1.c();
    					if_block1.m(td, t1);
    				}
    			}

    			if (!current || changed.columns && td_class_value !== (td_class_value = "relative " + ctx.column.class + " svelte-13j4zi7")) {
    				attr_dev(td, "class", td_class_value);
    			}

    			if (changed.columns || changed.editable || changed.columns) {
    				toggle_class(td, "cursor-pointer", ctx.editable && ctx.column.editable !== false);
    			}
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(if_block0);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(if_block0);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(td);
    			if (if_block0) if_block0.d();
    			if_block1.d();
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_each_block_1$1.name,
    		type: "each",
    		source: "(164:10) {#each columns as column, i}",
    		ctx
    	});

    	return block;
    }

    // (155:4) {#each sorted as item, j}
    function create_each_block$2(ctx) {
    	let tr;
    	let t;
    	let dispose_item_slot;
    	let current;
    	const item_slot_template = ctx.$$slots.item;
    	const item_slot = create_slot(item_slot_template, ctx, get_item_slot_context$1);
    	let each_value_1 = ctx.columns;
    	let each_blocks = [];

    	for (let i = 0; i < each_value_1.length; i += 1) {
    		each_blocks[i] = create_each_block_1$1(get_each_context_1$1(ctx, each_value_1, i));
    	}

    	const out = i => transition_out(each_blocks[i], 1, 1, () => {
    		each_blocks[i] = null;
    	});

    	function click_handler_1(...args) {
    		return ctx.click_handler_1(ctx, ...args);
    	}

    	const block = {
    		c: function create() {
    			if (!item_slot) {
    				tr = element("tr");

    				for (let i = 0; i < each_blocks.length; i += 1) {
    					each_blocks[i].c();
    				}

    				t = space();
    			}

    			if (item_slot) item_slot.c();

    			if (!item_slot) {
    				attr_dev(tr, "class", "svelte-13j4zi7");
    				toggle_class(tr, "selected", ctx.editing[ctx.j]);
    				add_location(tr, file$b, 156, 8, 5506);
    				dispose_item_slot = listen_dev(tr, "click", click_handler_1, false, false, false);
    			}
    		},
    		m: function mount(target, anchor) {
    			if (!item_slot) {
    				insert_dev(target, tr, anchor);

    				for (let i = 0; i < each_blocks.length; i += 1) {
    					each_blocks[i].m(tr, null);
    				}

    				insert_dev(target, t, anchor);
    			}

    			if (item_slot) {
    				item_slot.m(target, anchor);
    			}

    			current = true;
    		},
    		p: function update(changed, new_ctx) {
    			ctx = new_ctx;

    			if (!item_slot) {
    				if (changed.columns || changed.editable || changed.sorted || changed.editing || changed.dispatch || changed.$$scope) {
    					each_value_1 = ctx.columns;
    					let i;

    					for (i = 0; i < each_value_1.length; i += 1) {
    						const child_ctx = get_each_context_1$1(ctx, each_value_1, i);

    						if (each_blocks[i]) {
    							each_blocks[i].p(changed, child_ctx);
    							transition_in(each_blocks[i], 1);
    						} else {
    							each_blocks[i] = create_each_block_1$1(child_ctx);
    							each_blocks[i].c();
    							transition_in(each_blocks[i], 1);
    							each_blocks[i].m(tr, null);
    						}
    					}

    					group_outros();

    					for (i = each_value_1.length; i < each_blocks.length; i += 1) {
    						out(i);
    					}

    					check_outros();
    				}

    				if (changed.editing) {
    					toggle_class(tr, "selected", ctx.editing[ctx.j]);
    				}
    			}

    			if (item_slot && item_slot.p && changed.$$scope) {
    				item_slot.p(get_slot_changes(item_slot_template, ctx, changed, get_item_slot_changes$1), get_slot_context(item_slot_template, ctx, get_item_slot_context$1));
    			}
    		},
    		i: function intro(local) {
    			if (current) return;

    			for (let i = 0; i < each_value_1.length; i += 1) {
    				transition_in(each_blocks[i]);
    			}

    			transition_in(item_slot, local);
    			current = true;
    		},
    		o: function outro(local) {
    			each_blocks = each_blocks.filter(Boolean);

    			for (let i = 0; i < each_blocks.length; i += 1) {
    				transition_out(each_blocks[i]);
    			}

    			transition_out(item_slot, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (!item_slot) {
    				if (detaching) detach_dev(tr);
    				destroy_each(each_blocks, detaching);
    				if (detaching) detach_dev(t);
    				dispose_item_slot();
    			}

    			if (item_slot) item_slot.d(detaching);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_each_block$2.name,
    		type: "each",
    		source: "(155:4) {#each sorted as item, j}",
    		ctx
    	});

    	return block;
    }

    // (200:2) {#if pagination}
    function create_if_block$6(ctx) {
    	let tfoot;
    	let tr;
    	let td;
    	let div2;
    	let t0;
    	let div0;
    	let t2;
    	let updating_value;
    	let t3;
    	let t4;
    	let div1;
    	let t5;
    	let t6;

    	let t7_value = (ctx.offset + ctx.perPage > ctx.data.length
    	? ctx.data.length
    	: ctx.offset + ctx.perPage) + "";

    	let t7;
    	let t8;
    	let t9_value = ctx.data.length + "";
    	let t9;
    	let t10;
    	let t11;
    	let current;
    	const pagination_slot_template = ctx.$$slots.pagination;
    	const pagination_slot = create_slot(pagination_slot_template, ctx, get_pagination_slot_context);
    	const spacer0 = new Spacer$1({ $$inline: true });

    	function select_value_binding(value) {
    		ctx.select_value_binding.call(null, value);
    	}

    	let select_props = {
    		class: "w-16 h-8 mb-5",
    		remove: "bg-gray-300 bg-gray-100 select",
    		replace: { "pt-6": "pt-4" },
    		inputWrapperClasses: func,
    		appendClasses: func_1,
    		noUnderline: true,
    		dense: true,
    		items: ctx.perPageOptions
    	};

    	if (ctx.perPage !== void 0) {
    		select_props.value = ctx.perPage;
    	}

    	const select = new Select({ props: select_props, $$inline: true });
    	binding_callbacks.push(() => bind(select, "value", select_value_binding));
    	const spacer1 = new Spacer$1({ $$inline: true });

    	const button0_spread_levels = [
    		{ disabled: ctx.page - 1 < 1 },
    		{ icon: "keyboard_arrow_left" },
    		ctx.paginatorProps
    	];

    	let button0_props = {};

    	for (let i = 0; i < button0_spread_levels.length; i += 1) {
    		button0_props = assign(button0_props, button0_spread_levels[i]);
    	}

    	const button0 = new Button({ props: button0_props, $$inline: true });
    	button0.$on("click", ctx.click_handler_2);

    	const button1_spread_levels = [
    		{ disabled: ctx.page === ctx.pagesCount },
    		{ icon: "keyboard_arrow_right" },
    		ctx.paginatorProps
    	];

    	let button1_props = {};

    	for (let i = 0; i < button1_spread_levels.length; i += 1) {
    		button1_props = assign(button1_props, button1_spread_levels[i]);
    	}

    	const button1 = new Button({ props: button1_props, $$inline: true });
    	button1.$on("click", ctx.click_handler_3);

    	const block = {
    		c: function create() {
    			if (!pagination_slot) {
    				tfoot = element("tfoot");
    				tr = element("tr");
    				td = element("td");
    				div2 = element("div");
    				create_component(spacer0.$$.fragment);
    				t0 = space();
    				div0 = element("div");
    				div0.textContent = "Rows per page:";
    				t2 = space();
    				create_component(select.$$.fragment);
    				t3 = space();
    				create_component(spacer1.$$.fragment);
    				t4 = space();
    				div1 = element("div");
    				t5 = text(ctx.offset);
    				t6 = text("-");
    				t7 = text(t7_value);
    				t8 = text(" of ");
    				t9 = text(t9_value);
    				t10 = space();
    				create_component(button0.$$.fragment);
    				t11 = space();
    				create_component(button1.$$.fragment);
    			}

    			if (pagination_slot) pagination_slot.c();

    			if (!pagination_slot) {
    				attr_dev(div0, "class", "mr-1 py-1");
    				add_location(div0, file$b, 206, 14, 7238);
    				add_location(div1, file$b, 221, 14, 7831);
    				attr_dev(div2, "class", "flex justify-between items-center text-gray-700 text-sm w-full h-8");
    				add_location(div2, file$b, 204, 12, 7118);
    				attr_dev(td, "colspan", "100%");
    				attr_dev(td, "class", "svelte-13j4zi7");
    				add_location(td, file$b, 203, 10, 7086);
    				attr_dev(tr, "class", "svelte-13j4zi7");
    				add_location(tr, file$b, 202, 8, 7071);
    				add_location(tfoot, file$b, 201, 6, 7055);
    			}
    		},
    		m: function mount(target, anchor) {
    			if (!pagination_slot) {
    				insert_dev(target, tfoot, anchor);
    				append_dev(tfoot, tr);
    				append_dev(tr, td);
    				append_dev(td, div2);
    				mount_component(spacer0, div2, null);
    				append_dev(div2, t0);
    				append_dev(div2, div0);
    				append_dev(div2, t2);
    				mount_component(select, div2, null);
    				append_dev(div2, t3);
    				mount_component(spacer1, div2, null);
    				append_dev(div2, t4);
    				append_dev(div2, div1);
    				append_dev(div1, t5);
    				append_dev(div1, t6);
    				append_dev(div1, t7);
    				append_dev(div1, t8);
    				append_dev(div1, t9);
    				append_dev(div2, t10);
    				mount_component(button0, div2, null);
    				append_dev(div2, t11);
    				mount_component(button1, div2, null);
    			}

    			if (pagination_slot) {
    				pagination_slot.m(target, anchor);
    			}

    			current = true;
    		},
    		p: function update(changed, ctx) {
    			if (!pagination_slot) {
    				const select_changes = {};
    				if (changed.perPageOptions) select_changes.items = ctx.perPageOptions;

    				if (!updating_value && changed.perPage) {
    					updating_value = true;
    					select_changes.value = ctx.perPage;
    					add_flush_callback(() => updating_value = false);
    				}

    				select.$set(select_changes);
    				if (!current || changed.offset) set_data_dev(t5, ctx.offset);

    				if ((!current || (changed.offset || changed.perPage || changed.data)) && t7_value !== (t7_value = (ctx.offset + ctx.perPage > ctx.data.length
    				? ctx.data.length
    				: ctx.offset + ctx.perPage) + "")) set_data_dev(t7, t7_value);

    				if ((!current || changed.data) && t9_value !== (t9_value = ctx.data.length + "")) set_data_dev(t9, t9_value);

    				const button0_changes = changed.page || changed.paginatorProps
    				? get_spread_update(button0_spread_levels, [
    						changed.page && ({ disabled: ctx.page - 1 < 1 }),
    						button0_spread_levels[1],
    						changed.paginatorProps && get_spread_object(ctx.paginatorProps)
    					])
    				: {};

    				button0.$set(button0_changes);

    				const button1_changes = changed.page || changed.pagesCount || changed.paginatorProps
    				? get_spread_update(button1_spread_levels, [
    						(changed.page || changed.pagesCount) && ({ disabled: ctx.page === ctx.pagesCount }),
    						button1_spread_levels[1],
    						changed.paginatorProps && get_spread_object(ctx.paginatorProps)
    					])
    				: {};

    				button1.$set(button1_changes);
    			}

    			if (pagination_slot && pagination_slot.p && changed.$$scope) {
    				pagination_slot.p(get_slot_changes(pagination_slot_template, ctx, changed, get_pagination_slot_changes), get_slot_context(pagination_slot_template, ctx, get_pagination_slot_context));
    			}
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(spacer0.$$.fragment, local);
    			transition_in(select.$$.fragment, local);
    			transition_in(spacer1.$$.fragment, local);
    			transition_in(button0.$$.fragment, local);
    			transition_in(button1.$$.fragment, local);
    			transition_in(pagination_slot, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(spacer0.$$.fragment, local);
    			transition_out(select.$$.fragment, local);
    			transition_out(spacer1.$$.fragment, local);
    			transition_out(button0.$$.fragment, local);
    			transition_out(button1.$$.fragment, local);
    			transition_out(pagination_slot, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (!pagination_slot) {
    				if (detaching) detach_dev(tfoot);
    				destroy_component(spacer0);
    				destroy_component(select);
    				destroy_component(spacer1);
    				destroy_component(button0);
    				destroy_component(button1);
    			}

    			if (pagination_slot) pagination_slot.d(detaching);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_if_block$6.name,
    		type: "if",
    		source: "(200:2) {#if pagination}",
    		ctx
    	});

    	return block;
    }

    function create_fragment$b(ctx) {
    	let table_1;
    	let thead;
    	let t0;
    	let t1;
    	let tbody;
    	let t2;
    	let t3;
    	let current;
    	let each_value_2 = ctx.columns;
    	let each_blocks_1 = [];

    	for (let i = 0; i < each_value_2.length; i += 1) {
    		each_blocks_1[i] = create_each_block_2(get_each_context_2(ctx, each_value_2, i));
    	}

    	const out = i => transition_out(each_blocks_1[i], 1, 1, () => {
    		each_blocks_1[i] = null;
    	});

    	let if_block0 = ctx.loading && !ctx.hideProgress && create_if_block_3$1(ctx);
    	let each_value = ctx.sorted;
    	let each_blocks = [];

    	for (let i = 0; i < each_value.length; i += 1) {
    		each_blocks[i] = create_each_block$2(get_each_context$2(ctx, each_value, i));
    	}

    	const out_1 = i => transition_out(each_blocks[i], 1, 1, () => {
    		each_blocks[i] = null;
    	});

    	let if_block1 = ctx.pagination && create_if_block$6(ctx);
    	const footer_slot_template = ctx.$$slots.footer;
    	const footer_slot = create_slot(footer_slot_template, ctx, get_footer_slot_context);

    	const block = {
    		c: function create() {
    			table_1 = element("table");
    			thead = element("thead");

    			for (let i = 0; i < each_blocks_1.length; i += 1) {
    				each_blocks_1[i].c();
    			}

    			t0 = space();
    			if (if_block0) if_block0.c();
    			t1 = space();
    			tbody = element("tbody");

    			for (let i = 0; i < each_blocks.length; i += 1) {
    				each_blocks[i].c();
    			}

    			t2 = space();
    			if (if_block1) if_block1.c();
    			t3 = space();
    			if (footer_slot) footer_slot.c();
    			attr_dev(thead, "class", "items-center");
    			add_location(thead, file$b, 121, 2, 4449);
    			add_location(tbody, file$b, 153, 2, 5435);
    			attr_dev(table_1, "class", ctx.wrapperClasses);
    			add_location(table_1, file$b, 120, 0, 4398);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, table_1, anchor);
    			append_dev(table_1, thead);

    			for (let i = 0; i < each_blocks_1.length; i += 1) {
    				each_blocks_1[i].m(thead, null);
    			}

    			append_dev(table_1, t0);
    			if (if_block0) if_block0.m(table_1, null);
    			append_dev(table_1, t1);
    			append_dev(table_1, tbody);

    			for (let i = 0; i < each_blocks.length; i += 1) {
    				each_blocks[i].m(tbody, null);
    			}

    			append_dev(table_1, t2);
    			if (if_block1) if_block1.m(table_1, null);
    			append_dev(table_1, t3);

    			if (footer_slot) {
    				footer_slot.m(table_1, null);
    			}

    			ctx.table_1_binding(table_1);
    			current = true;
    		},
    		p: function update(changed, ctx) {
    			if (changed.sortable || changed.columns || changed.dispatch || changed.editing || changed.asc || changed.sortBy || changed.$$scope) {
    				each_value_2 = ctx.columns;
    				let i;

    				for (i = 0; i < each_value_2.length; i += 1) {
    					const child_ctx = get_each_context_2(ctx, each_value_2, i);

    					if (each_blocks_1[i]) {
    						each_blocks_1[i].p(changed, child_ctx);
    						transition_in(each_blocks_1[i], 1);
    					} else {
    						each_blocks_1[i] = create_each_block_2(child_ctx);
    						each_blocks_1[i].c();
    						transition_in(each_blocks_1[i], 1);
    						each_blocks_1[i].m(thead, null);
    					}
    				}

    				group_outros();

    				for (i = each_value_2.length; i < each_blocks_1.length; i += 1) {
    					out(i);
    				}

    				check_outros();
    			}

    			if (ctx.loading && !ctx.hideProgress) {
    				if (!if_block0) {
    					if_block0 = create_if_block_3$1(ctx);
    					if_block0.c();
    					transition_in(if_block0, 1);
    					if_block0.m(table_1, t1);
    				} else {
    					transition_in(if_block0, 1);
    				}
    			} else if (if_block0) {
    				group_outros();

    				transition_out(if_block0, 1, 1, () => {
    					if_block0 = null;
    				});

    				check_outros();
    			}

    			if (changed.editing || changed.editable || changed.columns || changed.sorted || changed.dispatch || changed.$$scope) {
    				each_value = ctx.sorted;
    				let i;

    				for (i = 0; i < each_value.length; i += 1) {
    					const child_ctx = get_each_context$2(ctx, each_value, i);

    					if (each_blocks[i]) {
    						each_blocks[i].p(changed, child_ctx);
    						transition_in(each_blocks[i], 1);
    					} else {
    						each_blocks[i] = create_each_block$2(child_ctx);
    						each_blocks[i].c();
    						transition_in(each_blocks[i], 1);
    						each_blocks[i].m(tbody, null);
    					}
    				}

    				group_outros();

    				for (i = each_value.length; i < each_blocks.length; i += 1) {
    					out_1(i);
    				}

    				check_outros();
    			}

    			if (ctx.pagination) {
    				if (if_block1) {
    					if_block1.p(changed, ctx);
    					transition_in(if_block1, 1);
    				} else {
    					if_block1 = create_if_block$6(ctx);
    					if_block1.c();
    					transition_in(if_block1, 1);
    					if_block1.m(table_1, t3);
    				}
    			} else if (if_block1) {
    				group_outros();

    				transition_out(if_block1, 1, 1, () => {
    					if_block1 = null;
    				});

    				check_outros();
    			}

    			if (footer_slot && footer_slot.p && changed.$$scope) {
    				footer_slot.p(get_slot_changes(footer_slot_template, ctx, changed, get_footer_slot_changes), get_slot_context(footer_slot_template, ctx, get_footer_slot_context));
    			}

    			if (!current || changed.wrapperClasses) {
    				attr_dev(table_1, "class", ctx.wrapperClasses);
    			}
    		},
    		i: function intro(local) {
    			if (current) return;

    			for (let i = 0; i < each_value_2.length; i += 1) {
    				transition_in(each_blocks_1[i]);
    			}

    			transition_in(if_block0);

    			for (let i = 0; i < each_value.length; i += 1) {
    				transition_in(each_blocks[i]);
    			}

    			transition_in(if_block1);
    			transition_in(footer_slot, local);
    			current = true;
    		},
    		o: function outro(local) {
    			each_blocks_1 = each_blocks_1.filter(Boolean);

    			for (let i = 0; i < each_blocks_1.length; i += 1) {
    				transition_out(each_blocks_1[i]);
    			}

    			transition_out(if_block0);
    			each_blocks = each_blocks.filter(Boolean);

    			for (let i = 0; i < each_blocks.length; i += 1) {
    				transition_out(each_blocks[i]);
    			}

    			transition_out(if_block1);
    			transition_out(footer_slot, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(table_1);
    			destroy_each(each_blocks_1, detaching);
    			if (if_block0) if_block0.d();
    			destroy_each(each_blocks, detaching);
    			if (if_block1) if_block1.d();
    			if (footer_slot) footer_slot.d(detaching);
    			ctx.table_1_binding(null);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$b.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    const func = c => c.replace("mt-2", "").replace("pb-6", "");
    const func_1 = c => c.replace("pt-4", "pt-3").replace("pr-4", "pr-2");

    function instance$8($$self, $$props, $$invalidate) {
    	let { data = [] } = $$props;

    	let { columns = Object.keys(data[0] || ({})).map(i => ({
    		label: (i || "").replace("_", " "),
    		field: i
    	})) } = $$props;

    	let { page = 1 } = $$props;
    	let { sort: sort$1 = sort } = $$props;
    	let { perPage = 10 } = $$props;
    	let { perPageOptions = [10, 20, 50] } = $$props;
    	let { asc = false } = $$props;
    	let { loading = false } = $$props;
    	let { hideProgress = false } = $$props;
    	let { wrapperClasses = "rounded elevation-3 relative text-sm overflow-x-auto" } = $$props;
    	let { editable = true } = $$props;
    	let { sortable = true } = $$props;
    	let { pagination = true } = $$props;

    	let { paginatorProps = {
    		color: "gray",
    		text: true,
    		flat: true,
    		dark: true,
    		remove: "px-4 px-3",
    		iconClasses: c => c.replace("p-4", ""),
    		disabledClasses: c => c.replace("text-white", "text-gray-200").replace("bg-gray-300", "bg-transparent").replace("text-gray-700", "")
    	} } = $$props;

    	let table = "";
    	let sortBy = null;
    	const dispatch = createEventDispatcher();
    	let editing = false;

    	const writable_props = [
    		"data",
    		"columns",
    		"page",
    		"sort",
    		"perPage",
    		"perPageOptions",
    		"asc",
    		"loading",
    		"hideProgress",
    		"wrapperClasses",
    		"editable",
    		"sortable",
    		"pagination",
    		"paginatorProps"
    	];

    	Object_1.keys($$props).forEach(key => {
    		if (!writable_props.includes(key) && !key.startsWith("$$")) console.warn(`<DataTable> was created with unknown prop '${key}'`);
    	});

    	let { $$slots = {}, $$scope } = $$props;

    	const click_handler = ({ column }) => {
    		if (column.sortable === false || !sortable) return;
    		dispatch("sort", column);
    		$$invalidate("editing", editing = false);
    		$$invalidate("asc", asc = sortBy === column ? !asc : false);
    		$$invalidate("sortBy", sortBy = column);
    	};

    	function change_handler(event) {
    		bubble($$self, event);
    	}

    	const blur_handler = ({ item, column }, { target }) => {
    		$$invalidate("editing", editing = false);
    		dispatch("update", { item, column, value: target.value });
    	};

    	const click_handler_1 = ({ j }, e) => {
    		if (!editable) return;

    		$$invalidate("editing", editing = {
    			[j]: (e.path.find(a => a.localName === "td") || ({})).cellIndex
    		});
    	};

    	function select_value_binding(value) {
    		perPage = value;
    		(($$invalidate("perPage", perPage), $$invalidate("pagination", pagination)), $$invalidate("data", data));
    	}

    	const click_handler_2 = () => {
    		$$invalidate("page", page -= 1);
    		table.scrollIntoView({ behavior: "smooth" });
    	};

    	const click_handler_3 = () => {
    		$$invalidate("page", page += 1);
    		table.scrollIntoView({ behavior: "smooth" });
    	};

    	function table_1_binding($$value) {
    		binding_callbacks[$$value ? "unshift" : "push"](() => {
    			$$invalidate("table", table = $$value);
    		});
    	}

    	$$self.$set = $$props => {
    		if ("data" in $$props) $$invalidate("data", data = $$props.data);
    		if ("columns" in $$props) $$invalidate("columns", columns = $$props.columns);
    		if ("page" in $$props) $$invalidate("page", page = $$props.page);
    		if ("sort" in $$props) $$invalidate("sort", sort$1 = $$props.sort);
    		if ("perPage" in $$props) $$invalidate("perPage", perPage = $$props.perPage);
    		if ("perPageOptions" in $$props) $$invalidate("perPageOptions", perPageOptions = $$props.perPageOptions);
    		if ("asc" in $$props) $$invalidate("asc", asc = $$props.asc);
    		if ("loading" in $$props) $$invalidate("loading", loading = $$props.loading);
    		if ("hideProgress" in $$props) $$invalidate("hideProgress", hideProgress = $$props.hideProgress);
    		if ("wrapperClasses" in $$props) $$invalidate("wrapperClasses", wrapperClasses = $$props.wrapperClasses);
    		if ("editable" in $$props) $$invalidate("editable", editable = $$props.editable);
    		if ("sortable" in $$props) $$invalidate("sortable", sortable = $$props.sortable);
    		if ("pagination" in $$props) $$invalidate("pagination", pagination = $$props.pagination);
    		if ("paginatorProps" in $$props) $$invalidate("paginatorProps", paginatorProps = $$props.paginatorProps);
    		if ("$$scope" in $$props) $$invalidate("$$scope", $$scope = $$props.$$scope);
    	};

    	$$self.$capture_state = () => {
    		return {
    			data,
    			columns,
    			page,
    			sort: sort$1,
    			perPage,
    			perPageOptions,
    			asc,
    			loading,
    			hideProgress,
    			wrapperClasses,
    			editable,
    			sortable,
    			pagination,
    			paginatorProps,
    			table,
    			sortBy,
    			editing,
    			offset,
    			sorted,
    			pagesCount
    		};
    	};

    	$$self.$inject_state = $$props => {
    		if ("data" in $$props) $$invalidate("data", data = $$props.data);
    		if ("columns" in $$props) $$invalidate("columns", columns = $$props.columns);
    		if ("page" in $$props) $$invalidate("page", page = $$props.page);
    		if ("sort" in $$props) $$invalidate("sort", sort$1 = $$props.sort);
    		if ("perPage" in $$props) $$invalidate("perPage", perPage = $$props.perPage);
    		if ("perPageOptions" in $$props) $$invalidate("perPageOptions", perPageOptions = $$props.perPageOptions);
    		if ("asc" in $$props) $$invalidate("asc", asc = $$props.asc);
    		if ("loading" in $$props) $$invalidate("loading", loading = $$props.loading);
    		if ("hideProgress" in $$props) $$invalidate("hideProgress", hideProgress = $$props.hideProgress);
    		if ("wrapperClasses" in $$props) $$invalidate("wrapperClasses", wrapperClasses = $$props.wrapperClasses);
    		if ("editable" in $$props) $$invalidate("editable", editable = $$props.editable);
    		if ("sortable" in $$props) $$invalidate("sortable", sortable = $$props.sortable);
    		if ("pagination" in $$props) $$invalidate("pagination", pagination = $$props.pagination);
    		if ("paginatorProps" in $$props) $$invalidate("paginatorProps", paginatorProps = $$props.paginatorProps);
    		if ("table" in $$props) $$invalidate("table", table = $$props.table);
    		if ("sortBy" in $$props) $$invalidate("sortBy", sortBy = $$props.sortBy);
    		if ("editing" in $$props) $$invalidate("editing", editing = $$props.editing);
    		if ("offset" in $$props) $$invalidate("offset", offset = $$props.offset);
    		if ("sorted" in $$props) $$invalidate("sorted", sorted = $$props.sorted);
    		if ("pagesCount" in $$props) $$invalidate("pagesCount", pagesCount = $$props.pagesCount);
    	};

    	let offset;
    	let sorted;
    	let pagesCount;

    	$$self.$$.update = (changed = { pagination: 1, perPage: 1, data: 1, page: 1, sort: 1, sortBy: 1, asc: 1, offset: 1 }) => {
    		if (changed.pagination || changed.perPage || changed.data) {
    			 {
    				$$invalidate("perPage", perPage = pagination ? perPage : data.length);
    				$$invalidate("page", page = 1);
    			}
    		}

    		if (changed.page || changed.perPage) {
    			 $$invalidate("offset", offset = page * perPage - perPage);
    		}

    		if (changed.sort || changed.data || changed.sortBy || changed.asc || changed.offset || changed.perPage) {
    			 $$invalidate("sorted", sorted = sort$1(data, sortBy, asc).slice(offset, perPage + offset));
    		}

    		if (changed.data || changed.perPage) {
    			 $$invalidate("pagesCount", pagesCount = Math.ceil(data.length / perPage));
    		}
    	};

    	return {
    		data,
    		columns,
    		page,
    		sort: sort$1,
    		perPage,
    		perPageOptions,
    		asc,
    		loading,
    		hideProgress,
    		wrapperClasses,
    		editable,
    		sortable,
    		pagination,
    		paginatorProps,
    		table,
    		sortBy,
    		dispatch,
    		editing,
    		offset,
    		sorted,
    		pagesCount,
    		click_handler,
    		change_handler,
    		blur_handler,
    		click_handler_1,
    		select_value_binding,
    		click_handler_2,
    		click_handler_3,
    		table_1_binding,
    		$$slots,
    		$$scope
    	};
    }

    class DataTable extends SvelteComponentDev {
    	constructor(options) {
    		super(options);

    		init(this, options, instance$8, create_fragment$b, safe_not_equal, {
    			data: 0,
    			columns: 0,
    			page: 0,
    			sort: 0,
    			perPage: 0,
    			perPageOptions: 0,
    			asc: 0,
    			loading: 0,
    			hideProgress: 0,
    			wrapperClasses: 0,
    			editable: 0,
    			sortable: 0,
    			pagination: 0,
    			paginatorProps: 0
    		});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "DataTable",
    			options,
    			id: create_fragment$b.name
    		});
    	}

    	get data() {
    		throw new Error("<DataTable>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set data(value) {
    		throw new Error("<DataTable>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get columns() {
    		throw new Error("<DataTable>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set columns(value) {
    		throw new Error("<DataTable>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get page() {
    		throw new Error("<DataTable>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set page(value) {
    		throw new Error("<DataTable>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get sort() {
    		throw new Error("<DataTable>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set sort(value) {
    		throw new Error("<DataTable>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get perPage() {
    		throw new Error("<DataTable>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set perPage(value) {
    		throw new Error("<DataTable>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get perPageOptions() {
    		throw new Error("<DataTable>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set perPageOptions(value) {
    		throw new Error("<DataTable>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get asc() {
    		throw new Error("<DataTable>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set asc(value) {
    		throw new Error("<DataTable>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get loading() {
    		throw new Error("<DataTable>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set loading(value) {
    		throw new Error("<DataTable>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get hideProgress() {
    		throw new Error("<DataTable>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set hideProgress(value) {
    		throw new Error("<DataTable>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get wrapperClasses() {
    		throw new Error("<DataTable>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set wrapperClasses(value) {
    		throw new Error("<DataTable>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get editable() {
    		throw new Error("<DataTable>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set editable(value) {
    		throw new Error("<DataTable>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get sortable() {
    		throw new Error("<DataTable>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set sortable(value) {
    		throw new Error("<DataTable>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get pagination() {
    		throw new Error("<DataTable>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set pagination(value) {
    		throw new Error("<DataTable>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	get paginatorProps() {
    		throw new Error("<DataTable>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set paginatorProps(value) {
    		throw new Error("<DataTable>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}
    }

    /* src/components/panels/PanelSend.svelte generated by Svelte v3.14.0 */
    const file$c = "src/components/panels/PanelSend.svelte";

    // (12:3) <Button id="dialogbtn" v-on:click.native="btnClick">
    function create_default_slot$5(ctx) {
    	let t;

    	const block = {
    		c: function create() {
    			t = text("Send");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, t, anchor);
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(t);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_default_slot$5.name,
    		type: "slot",
    		source: "(12:3) <Button id=\\\"dialogbtn\\\" v-on:click.native=\\\"btnClick\\\">",
    		ctx
    	});

    	return block;
    }

    function create_fragment$c(ctx) {
    	let div2;
    	let div1;
    	let t0;
    	let div0;
    	let t1;
    	let current;

    	const textfield0 = new TextField({
    			props: {
    				id: "sendAddress",
    				label: "Test label",
    				placeholder: "Enter DUO address",
    				outlined: true,
    				hint: "Test hint",
    				class: "e-outline flx noMargin fii bgfff"
    			},
    			$$inline: true
    		});

    	const textfield1 = new TextField({
    			props: {
    				label: "Test label",
    				class: "e-outline noMargin fii bgfff"
    			},
    			$$inline: true
    		});

    	const button = new Button({
    			props: {
    				id: "dialogbtn",
    				"v-on:click.native": "btnClick",
    				$$slots: { default: [create_default_slot$5] },
    				$$scope: { ctx }
    			},
    			$$inline: true
    		});

    	const block = {
    		c: function create() {
    			div2 = element("div");
    			div1 = element("div");
    			create_component(textfield0.$$.fragment);
    			t0 = space();
    			div0 = element("div");
    			create_component(textfield1.$$.fragment);
    			t1 = space();
    			create_component(button.$$.fragment);
    			attr_dev(div0, "class", "flx fii");
    			add_location(div0, file$c, 9, 3, 314);
    			attr_dev(div1, "class", "flx flc fii justifyBetween");
    			add_location(div1, file$c, 7, 1, 119);
    			attr_dev(div2, "class", "rwrap flx");
    			add_location(div2, file$c, 6, 0, 94);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div2, anchor);
    			append_dev(div2, div1);
    			mount_component(textfield0, div1, null);
    			append_dev(div1, t0);
    			append_dev(div1, div0);
    			mount_component(textfield1, div0, null);
    			append_dev(div0, t1);
    			mount_component(button, div0, null);
    			current = true;
    		},
    		p: function update(changed, ctx) {
    			const button_changes = {};

    			if (changed.$$scope) {
    				button_changes.$$scope = { changed, ctx };
    			}

    			button.$set(button_changes);
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(textfield0.$$.fragment, local);
    			transition_in(textfield1.$$.fragment, local);
    			transition_in(button.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(textfield0.$$.fragment, local);
    			transition_out(textfield1.$$.fragment, local);
    			transition_out(button.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div2);
    			destroy_component(textfield0);
    			destroy_component(textfield1);
    			destroy_component(button);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$c.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class PanelSend extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$c, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PanelSend",
    			options,
    			id: create_fragment$c.name
    		});
    	}
    }

    function ascending(a, b) {
      return a < b ? -1 : a > b ? 1 : a >= b ? 0 : NaN;
    }

    function bisector(compare) {
      if (compare.length === 1) compare = ascendingComparator(compare);
      return {
        left: function(a, x, lo, hi) {
          if (lo == null) lo = 0;
          if (hi == null) hi = a.length;
          while (lo < hi) {
            var mid = lo + hi >>> 1;
            if (compare(a[mid], x) < 0) lo = mid + 1;
            else hi = mid;
          }
          return lo;
        },
        right: function(a, x, lo, hi) {
          if (lo == null) lo = 0;
          if (hi == null) hi = a.length;
          while (lo < hi) {
            var mid = lo + hi >>> 1;
            if (compare(a[mid], x) > 0) hi = mid;
            else lo = mid + 1;
          }
          return lo;
        }
      };
    }

    function ascendingComparator(f) {
      return function(d, x) {
        return ascending(f(d), x);
      };
    }

    var ascendingBisect = bisector(ascending);
    var bisectRight = ascendingBisect.right;

    var e10 = Math.sqrt(50),
        e5 = Math.sqrt(10),
        e2 = Math.sqrt(2);

    function ticks(start, stop, count) {
      var reverse,
          i = -1,
          n,
          ticks,
          step;

      stop = +stop, start = +start, count = +count;
      if (start === stop && count > 0) return [start];
      if (reverse = stop < start) n = start, start = stop, stop = n;
      if ((step = tickIncrement(start, stop, count)) === 0 || !isFinite(step)) return [];

      if (step > 0) {
        start = Math.ceil(start / step);
        stop = Math.floor(stop / step);
        ticks = new Array(n = Math.ceil(stop - start + 1));
        while (++i < n) ticks[i] = (start + i) * step;
      } else {
        start = Math.floor(start * step);
        stop = Math.ceil(stop * step);
        ticks = new Array(n = Math.ceil(start - stop + 1));
        while (++i < n) ticks[i] = (start - i) / step;
      }

      if (reverse) ticks.reverse();

      return ticks;
    }

    function tickIncrement(start, stop, count) {
      var step = (stop - start) / Math.max(0, count),
          power = Math.floor(Math.log(step) / Math.LN10),
          error = step / Math.pow(10, power);
      return power >= 0
          ? (error >= e10 ? 10 : error >= e5 ? 5 : error >= e2 ? 2 : 1) * Math.pow(10, power)
          : -Math.pow(10, -power) / (error >= e10 ? 10 : error >= e5 ? 5 : error >= e2 ? 2 : 1);
    }

    function tickStep(start, stop, count) {
      var step0 = Math.abs(stop - start) / Math.max(0, count),
          step1 = Math.pow(10, Math.floor(Math.log(step0) / Math.LN10)),
          error = step0 / step1;
      if (error >= e10) step1 *= 10;
      else if (error >= e5) step1 *= 5;
      else if (error >= e2) step1 *= 2;
      return stop < start ? -step1 : step1;
    }

    function initRange(domain, range) {
      switch (arguments.length) {
        case 0: break;
        case 1: this.range(domain); break;
        default: this.range(range).domain(domain); break;
      }
      return this;
    }

    function define(constructor, factory, prototype) {
      constructor.prototype = factory.prototype = prototype;
      prototype.constructor = constructor;
    }

    function extend(parent, definition) {
      var prototype = Object.create(parent.prototype);
      for (var key in definition) prototype[key] = definition[key];
      return prototype;
    }

    function Color() {}

    var darker = 0.7;
    var brighter = 1 / darker;

    var reI = "\\s*([+-]?\\d+)\\s*",
        reN = "\\s*([+-]?\\d*\\.?\\d+(?:[eE][+-]?\\d+)?)\\s*",
        reP = "\\s*([+-]?\\d*\\.?\\d+(?:[eE][+-]?\\d+)?)%\\s*",
        reHex = /^#([0-9a-f]{3,8})$/,
        reRgbInteger = new RegExp("^rgb\\(" + [reI, reI, reI] + "\\)$"),
        reRgbPercent = new RegExp("^rgb\\(" + [reP, reP, reP] + "\\)$"),
        reRgbaInteger = new RegExp("^rgba\\(" + [reI, reI, reI, reN] + "\\)$"),
        reRgbaPercent = new RegExp("^rgba\\(" + [reP, reP, reP, reN] + "\\)$"),
        reHslPercent = new RegExp("^hsl\\(" + [reN, reP, reP] + "\\)$"),
        reHslaPercent = new RegExp("^hsla\\(" + [reN, reP, reP, reN] + "\\)$");

    var named = {
      aliceblue: 0xf0f8ff,
      antiquewhite: 0xfaebd7,
      aqua: 0x00ffff,
      aquamarine: 0x7fffd4,
      azure: 0xf0ffff,
      beige: 0xf5f5dc,
      bisque: 0xffe4c4,
      black: 0x000000,
      blanchedalmond: 0xffebcd,
      blue: 0x0000ff,
      blueviolet: 0x8a2be2,
      brown: 0xa52a2a,
      burlywood: 0xdeb887,
      cadetblue: 0x5f9ea0,
      chartreuse: 0x7fff00,
      chocolate: 0xd2691e,
      coral: 0xff7f50,
      cornflowerblue: 0x6495ed,
      cornsilk: 0xfff8dc,
      crimson: 0xdc143c,
      cyan: 0x00ffff,
      darkblue: 0x00008b,
      darkcyan: 0x008b8b,
      darkgoldenrod: 0xb8860b,
      darkgray: 0xa9a9a9,
      darkgreen: 0x006400,
      darkgrey: 0xa9a9a9,
      darkkhaki: 0xbdb76b,
      darkmagenta: 0x8b008b,
      darkolivegreen: 0x556b2f,
      darkorange: 0xff8c00,
      darkorchid: 0x9932cc,
      darkred: 0x8b0000,
      darksalmon: 0xe9967a,
      darkseagreen: 0x8fbc8f,
      darkslateblue: 0x483d8b,
      darkslategray: 0x2f4f4f,
      darkslategrey: 0x2f4f4f,
      darkturquoise: 0x00ced1,
      darkviolet: 0x9400d3,
      deeppink: 0xff1493,
      deepskyblue: 0x00bfff,
      dimgray: 0x696969,
      dimgrey: 0x696969,
      dodgerblue: 0x1e90ff,
      firebrick: 0xb22222,
      floralwhite: 0xfffaf0,
      forestgreen: 0x228b22,
      fuchsia: 0xff00ff,
      gainsboro: 0xdcdcdc,
      ghostwhite: 0xf8f8ff,
      gold: 0xffd700,
      goldenrod: 0xdaa520,
      gray: 0x808080,
      green: 0x008000,
      greenyellow: 0xadff2f,
      grey: 0x808080,
      honeydew: 0xf0fff0,
      hotpink: 0xff69b4,
      indianred: 0xcd5c5c,
      indigo: 0x4b0082,
      ivory: 0xfffff0,
      khaki: 0xf0e68c,
      lavender: 0xe6e6fa,
      lavenderblush: 0xfff0f5,
      lawngreen: 0x7cfc00,
      lemonchiffon: 0xfffacd,
      lightblue: 0xadd8e6,
      lightcoral: 0xf08080,
      lightcyan: 0xe0ffff,
      lightgoldenrodyellow: 0xfafad2,
      lightgray: 0xd3d3d3,
      lightgreen: 0x90ee90,
      lightgrey: 0xd3d3d3,
      lightpink: 0xffb6c1,
      lightsalmon: 0xffa07a,
      lightseagreen: 0x20b2aa,
      lightskyblue: 0x87cefa,
      lightslategray: 0x778899,
      lightslategrey: 0x778899,
      lightsteelblue: 0xb0c4de,
      lightyellow: 0xffffe0,
      lime: 0x00ff00,
      limegreen: 0x32cd32,
      linen: 0xfaf0e6,
      magenta: 0xff00ff,
      maroon: 0x800000,
      mediumaquamarine: 0x66cdaa,
      mediumblue: 0x0000cd,
      mediumorchid: 0xba55d3,
      mediumpurple: 0x9370db,
      mediumseagreen: 0x3cb371,
      mediumslateblue: 0x7b68ee,
      mediumspringgreen: 0x00fa9a,
      mediumturquoise: 0x48d1cc,
      mediumvioletred: 0xc71585,
      midnightblue: 0x191970,
      mintcream: 0xf5fffa,
      mistyrose: 0xffe4e1,
      moccasin: 0xffe4b5,
      navajowhite: 0xffdead,
      navy: 0x000080,
      oldlace: 0xfdf5e6,
      olive: 0x808000,
      olivedrab: 0x6b8e23,
      orange: 0xffa500,
      orangered: 0xff4500,
      orchid: 0xda70d6,
      palegoldenrod: 0xeee8aa,
      palegreen: 0x98fb98,
      paleturquoise: 0xafeeee,
      palevioletred: 0xdb7093,
      papayawhip: 0xffefd5,
      peachpuff: 0xffdab9,
      peru: 0xcd853f,
      pink: 0xffc0cb,
      plum: 0xdda0dd,
      powderblue: 0xb0e0e6,
      purple: 0x800080,
      rebeccapurple: 0x663399,
      red: 0xff0000,
      rosybrown: 0xbc8f8f,
      royalblue: 0x4169e1,
      saddlebrown: 0x8b4513,
      salmon: 0xfa8072,
      sandybrown: 0xf4a460,
      seagreen: 0x2e8b57,
      seashell: 0xfff5ee,
      sienna: 0xa0522d,
      silver: 0xc0c0c0,
      skyblue: 0x87ceeb,
      slateblue: 0x6a5acd,
      slategray: 0x708090,
      slategrey: 0x708090,
      snow: 0xfffafa,
      springgreen: 0x00ff7f,
      steelblue: 0x4682b4,
      tan: 0xd2b48c,
      teal: 0x008080,
      thistle: 0xd8bfd8,
      tomato: 0xff6347,
      turquoise: 0x40e0d0,
      violet: 0xee82ee,
      wheat: 0xf5deb3,
      white: 0xffffff,
      whitesmoke: 0xf5f5f5,
      yellow: 0xffff00,
      yellowgreen: 0x9acd32
    };

    define(Color, color, {
      copy: function(channels) {
        return Object.assign(new this.constructor, this, channels);
      },
      displayable: function() {
        return this.rgb().displayable();
      },
      hex: color_formatHex, // Deprecated! Use color.formatHex.
      formatHex: color_formatHex,
      formatHsl: color_formatHsl,
      formatRgb: color_formatRgb,
      toString: color_formatRgb
    });

    function color_formatHex() {
      return this.rgb().formatHex();
    }

    function color_formatHsl() {
      return hslConvert(this).formatHsl();
    }

    function color_formatRgb() {
      return this.rgb().formatRgb();
    }

    function color(format) {
      var m, l;
      format = (format + "").trim().toLowerCase();
      return (m = reHex.exec(format)) ? (l = m[1].length, m = parseInt(m[1], 16), l === 6 ? rgbn(m) // #ff0000
          : l === 3 ? new Rgb((m >> 8 & 0xf) | (m >> 4 & 0xf0), (m >> 4 & 0xf) | (m & 0xf0), ((m & 0xf) << 4) | (m & 0xf), 1) // #f00
          : l === 8 ? new Rgb(m >> 24 & 0xff, m >> 16 & 0xff, m >> 8 & 0xff, (m & 0xff) / 0xff) // #ff000000
          : l === 4 ? new Rgb((m >> 12 & 0xf) | (m >> 8 & 0xf0), (m >> 8 & 0xf) | (m >> 4 & 0xf0), (m >> 4 & 0xf) | (m & 0xf0), (((m & 0xf) << 4) | (m & 0xf)) / 0xff) // #f000
          : null) // invalid hex
          : (m = reRgbInteger.exec(format)) ? new Rgb(m[1], m[2], m[3], 1) // rgb(255, 0, 0)
          : (m = reRgbPercent.exec(format)) ? new Rgb(m[1] * 255 / 100, m[2] * 255 / 100, m[3] * 255 / 100, 1) // rgb(100%, 0%, 0%)
          : (m = reRgbaInteger.exec(format)) ? rgba(m[1], m[2], m[3], m[4]) // rgba(255, 0, 0, 1)
          : (m = reRgbaPercent.exec(format)) ? rgba(m[1] * 255 / 100, m[2] * 255 / 100, m[3] * 255 / 100, m[4]) // rgb(100%, 0%, 0%, 1)
          : (m = reHslPercent.exec(format)) ? hsla(m[1], m[2] / 100, m[3] / 100, 1) // hsl(120, 50%, 50%)
          : (m = reHslaPercent.exec(format)) ? hsla(m[1], m[2] / 100, m[3] / 100, m[4]) // hsla(120, 50%, 50%, 1)
          : named.hasOwnProperty(format) ? rgbn(named[format]) // eslint-disable-line no-prototype-builtins
          : format === "transparent" ? new Rgb(NaN, NaN, NaN, 0)
          : null;
    }

    function rgbn(n) {
      return new Rgb(n >> 16 & 0xff, n >> 8 & 0xff, n & 0xff, 1);
    }

    function rgba(r, g, b, a) {
      if (a <= 0) r = g = b = NaN;
      return new Rgb(r, g, b, a);
    }

    function rgbConvert(o) {
      if (!(o instanceof Color)) o = color(o);
      if (!o) return new Rgb;
      o = o.rgb();
      return new Rgb(o.r, o.g, o.b, o.opacity);
    }

    function rgb(r, g, b, opacity) {
      return arguments.length === 1 ? rgbConvert(r) : new Rgb(r, g, b, opacity == null ? 1 : opacity);
    }

    function Rgb(r, g, b, opacity) {
      this.r = +r;
      this.g = +g;
      this.b = +b;
      this.opacity = +opacity;
    }

    define(Rgb, rgb, extend(Color, {
      brighter: function(k) {
        k = k == null ? brighter : Math.pow(brighter, k);
        return new Rgb(this.r * k, this.g * k, this.b * k, this.opacity);
      },
      darker: function(k) {
        k = k == null ? darker : Math.pow(darker, k);
        return new Rgb(this.r * k, this.g * k, this.b * k, this.opacity);
      },
      rgb: function() {
        return this;
      },
      displayable: function() {
        return (-0.5 <= this.r && this.r < 255.5)
            && (-0.5 <= this.g && this.g < 255.5)
            && (-0.5 <= this.b && this.b < 255.5)
            && (0 <= this.opacity && this.opacity <= 1);
      },
      hex: rgb_formatHex, // Deprecated! Use color.formatHex.
      formatHex: rgb_formatHex,
      formatRgb: rgb_formatRgb,
      toString: rgb_formatRgb
    }));

    function rgb_formatHex() {
      return "#" + hex(this.r) + hex(this.g) + hex(this.b);
    }

    function rgb_formatRgb() {
      var a = this.opacity; a = isNaN(a) ? 1 : Math.max(0, Math.min(1, a));
      return (a === 1 ? "rgb(" : "rgba(")
          + Math.max(0, Math.min(255, Math.round(this.r) || 0)) + ", "
          + Math.max(0, Math.min(255, Math.round(this.g) || 0)) + ", "
          + Math.max(0, Math.min(255, Math.round(this.b) || 0))
          + (a === 1 ? ")" : ", " + a + ")");
    }

    function hex(value) {
      value = Math.max(0, Math.min(255, Math.round(value) || 0));
      return (value < 16 ? "0" : "") + value.toString(16);
    }

    function hsla(h, s, l, a) {
      if (a <= 0) h = s = l = NaN;
      else if (l <= 0 || l >= 1) h = s = NaN;
      else if (s <= 0) h = NaN;
      return new Hsl(h, s, l, a);
    }

    function hslConvert(o) {
      if (o instanceof Hsl) return new Hsl(o.h, o.s, o.l, o.opacity);
      if (!(o instanceof Color)) o = color(o);
      if (!o) return new Hsl;
      if (o instanceof Hsl) return o;
      o = o.rgb();
      var r = o.r / 255,
          g = o.g / 255,
          b = o.b / 255,
          min = Math.min(r, g, b),
          max = Math.max(r, g, b),
          h = NaN,
          s = max - min,
          l = (max + min) / 2;
      if (s) {
        if (r === max) h = (g - b) / s + (g < b) * 6;
        else if (g === max) h = (b - r) / s + 2;
        else h = (r - g) / s + 4;
        s /= l < 0.5 ? max + min : 2 - max - min;
        h *= 60;
      } else {
        s = l > 0 && l < 1 ? 0 : h;
      }
      return new Hsl(h, s, l, o.opacity);
    }

    function hsl(h, s, l, opacity) {
      return arguments.length === 1 ? hslConvert(h) : new Hsl(h, s, l, opacity == null ? 1 : opacity);
    }

    function Hsl(h, s, l, opacity) {
      this.h = +h;
      this.s = +s;
      this.l = +l;
      this.opacity = +opacity;
    }

    define(Hsl, hsl, extend(Color, {
      brighter: function(k) {
        k = k == null ? brighter : Math.pow(brighter, k);
        return new Hsl(this.h, this.s, this.l * k, this.opacity);
      },
      darker: function(k) {
        k = k == null ? darker : Math.pow(darker, k);
        return new Hsl(this.h, this.s, this.l * k, this.opacity);
      },
      rgb: function() {
        var h = this.h % 360 + (this.h < 0) * 360,
            s = isNaN(h) || isNaN(this.s) ? 0 : this.s,
            l = this.l,
            m2 = l + (l < 0.5 ? l : 1 - l) * s,
            m1 = 2 * l - m2;
        return new Rgb(
          hsl2rgb(h >= 240 ? h - 240 : h + 120, m1, m2),
          hsl2rgb(h, m1, m2),
          hsl2rgb(h < 120 ? h + 240 : h - 120, m1, m2),
          this.opacity
        );
      },
      displayable: function() {
        return (0 <= this.s && this.s <= 1 || isNaN(this.s))
            && (0 <= this.l && this.l <= 1)
            && (0 <= this.opacity && this.opacity <= 1);
      },
      formatHsl: function() {
        var a = this.opacity; a = isNaN(a) ? 1 : Math.max(0, Math.min(1, a));
        return (a === 1 ? "hsl(" : "hsla(")
            + (this.h || 0) + ", "
            + (this.s || 0) * 100 + "%, "
            + (this.l || 0) * 100 + "%"
            + (a === 1 ? ")" : ", " + a + ")");
      }
    }));

    /* From FvD 13.37, CSS Color Module Level 3 */
    function hsl2rgb(h, m1, m2) {
      return (h < 60 ? m1 + (m2 - m1) * h / 60
          : h < 180 ? m2
          : h < 240 ? m1 + (m2 - m1) * (240 - h) / 60
          : m1) * 255;
    }

    var deg2rad = Math.PI / 180;
    var rad2deg = 180 / Math.PI;

    // https://observablehq.com/@mbostock/lab-and-rgb
    var K = 18,
        Xn = 0.96422,
        Yn = 1,
        Zn = 0.82521,
        t0 = 4 / 29,
        t1 = 6 / 29,
        t2 = 3 * t1 * t1,
        t3 = t1 * t1 * t1;

    function labConvert(o) {
      if (o instanceof Lab) return new Lab(o.l, o.a, o.b, o.opacity);
      if (o instanceof Hcl) return hcl2lab(o);
      if (!(o instanceof Rgb)) o = rgbConvert(o);
      var r = rgb2lrgb(o.r),
          g = rgb2lrgb(o.g),
          b = rgb2lrgb(o.b),
          y = xyz2lab((0.2225045 * r + 0.7168786 * g + 0.0606169 * b) / Yn), x, z;
      if (r === g && g === b) x = z = y; else {
        x = xyz2lab((0.4360747 * r + 0.3850649 * g + 0.1430804 * b) / Xn);
        z = xyz2lab((0.0139322 * r + 0.0971045 * g + 0.7141733 * b) / Zn);
      }
      return new Lab(116 * y - 16, 500 * (x - y), 200 * (y - z), o.opacity);
    }

    function lab(l, a, b, opacity) {
      return arguments.length === 1 ? labConvert(l) : new Lab(l, a, b, opacity == null ? 1 : opacity);
    }

    function Lab(l, a, b, opacity) {
      this.l = +l;
      this.a = +a;
      this.b = +b;
      this.opacity = +opacity;
    }

    define(Lab, lab, extend(Color, {
      brighter: function(k) {
        return new Lab(this.l + K * (k == null ? 1 : k), this.a, this.b, this.opacity);
      },
      darker: function(k) {
        return new Lab(this.l - K * (k == null ? 1 : k), this.a, this.b, this.opacity);
      },
      rgb: function() {
        var y = (this.l + 16) / 116,
            x = isNaN(this.a) ? y : y + this.a / 500,
            z = isNaN(this.b) ? y : y - this.b / 200;
        x = Xn * lab2xyz(x);
        y = Yn * lab2xyz(y);
        z = Zn * lab2xyz(z);
        return new Rgb(
          lrgb2rgb( 3.1338561 * x - 1.6168667 * y - 0.4906146 * z),
          lrgb2rgb(-0.9787684 * x + 1.9161415 * y + 0.0334540 * z),
          lrgb2rgb( 0.0719453 * x - 0.2289914 * y + 1.4052427 * z),
          this.opacity
        );
      }
    }));

    function xyz2lab(t) {
      return t > t3 ? Math.pow(t, 1 / 3) : t / t2 + t0;
    }

    function lab2xyz(t) {
      return t > t1 ? t * t * t : t2 * (t - t0);
    }

    function lrgb2rgb(x) {
      return 255 * (x <= 0.0031308 ? 12.92 * x : 1.055 * Math.pow(x, 1 / 2.4) - 0.055);
    }

    function rgb2lrgb(x) {
      return (x /= 255) <= 0.04045 ? x / 12.92 : Math.pow((x + 0.055) / 1.055, 2.4);
    }

    function hclConvert(o) {
      if (o instanceof Hcl) return new Hcl(o.h, o.c, o.l, o.opacity);
      if (!(o instanceof Lab)) o = labConvert(o);
      if (o.a === 0 && o.b === 0) return new Hcl(NaN, 0 < o.l && o.l < 100 ? 0 : NaN, o.l, o.opacity);
      var h = Math.atan2(o.b, o.a) * rad2deg;
      return new Hcl(h < 0 ? h + 360 : h, Math.sqrt(o.a * o.a + o.b * o.b), o.l, o.opacity);
    }

    function hcl(h, c, l, opacity) {
      return arguments.length === 1 ? hclConvert(h) : new Hcl(h, c, l, opacity == null ? 1 : opacity);
    }

    function Hcl(h, c, l, opacity) {
      this.h = +h;
      this.c = +c;
      this.l = +l;
      this.opacity = +opacity;
    }

    function hcl2lab(o) {
      if (isNaN(o.h)) return new Lab(o.l, 0, 0, o.opacity);
      var h = o.h * deg2rad;
      return new Lab(o.l, Math.cos(h) * o.c, Math.sin(h) * o.c, o.opacity);
    }

    define(Hcl, hcl, extend(Color, {
      brighter: function(k) {
        return new Hcl(this.h, this.c, this.l + K * (k == null ? 1 : k), this.opacity);
      },
      darker: function(k) {
        return new Hcl(this.h, this.c, this.l - K * (k == null ? 1 : k), this.opacity);
      },
      rgb: function() {
        return hcl2lab(this).rgb();
      }
    }));

    var A = -0.14861,
        B = +1.78277,
        C = -0.29227,
        D = -0.90649,
        E = +1.97294,
        ED = E * D,
        EB = E * B,
        BC_DA = B * C - D * A;

    function cubehelixConvert(o) {
      if (o instanceof Cubehelix) return new Cubehelix(o.h, o.s, o.l, o.opacity);
      if (!(o instanceof Rgb)) o = rgbConvert(o);
      var r = o.r / 255,
          g = o.g / 255,
          b = o.b / 255,
          l = (BC_DA * b + ED * r - EB * g) / (BC_DA + ED - EB),
          bl = b - l,
          k = (E * (g - l) - C * bl) / D,
          s = Math.sqrt(k * k + bl * bl) / (E * l * (1 - l)), // NaN if l=0 or l=1
          h = s ? Math.atan2(k, bl) * rad2deg - 120 : NaN;
      return new Cubehelix(h < 0 ? h + 360 : h, s, l, o.opacity);
    }

    function cubehelix(h, s, l, opacity) {
      return arguments.length === 1 ? cubehelixConvert(h) : new Cubehelix(h, s, l, opacity == null ? 1 : opacity);
    }

    function Cubehelix(h, s, l, opacity) {
      this.h = +h;
      this.s = +s;
      this.l = +l;
      this.opacity = +opacity;
    }

    define(Cubehelix, cubehelix, extend(Color, {
      brighter: function(k) {
        k = k == null ? brighter : Math.pow(brighter, k);
        return new Cubehelix(this.h, this.s, this.l * k, this.opacity);
      },
      darker: function(k) {
        k = k == null ? darker : Math.pow(darker, k);
        return new Cubehelix(this.h, this.s, this.l * k, this.opacity);
      },
      rgb: function() {
        var h = isNaN(this.h) ? 0 : (this.h + 120) * deg2rad,
            l = +this.l,
            a = isNaN(this.s) ? 0 : this.s * l * (1 - l),
            cosh = Math.cos(h),
            sinh = Math.sin(h);
        return new Rgb(
          255 * (l + a * (A * cosh + B * sinh)),
          255 * (l + a * (C * cosh + D * sinh)),
          255 * (l + a * (E * cosh)),
          this.opacity
        );
      }
    }));

    function constant(x) {
      return function() {
        return x;
      };
    }

    function linear(a, d) {
      return function(t) {
        return a + t * d;
      };
    }

    function exponential(a, b, y) {
      return a = Math.pow(a, y), b = Math.pow(b, y) - a, y = 1 / y, function(t) {
        return Math.pow(a + t * b, y);
      };
    }

    function gamma(y) {
      return (y = +y) === 1 ? nogamma : function(a, b) {
        return b - a ? exponential(a, b, y) : constant(isNaN(a) ? b : a);
      };
    }

    function nogamma(a, b) {
      var d = b - a;
      return d ? linear(a, d) : constant(isNaN(a) ? b : a);
    }

    var rgb$1 = (function rgbGamma(y) {
      var color = gamma(y);

      function rgb$1(start, end) {
        var r = color((start = rgb(start)).r, (end = rgb(end)).r),
            g = color(start.g, end.g),
            b = color(start.b, end.b),
            opacity = nogamma(start.opacity, end.opacity);
        return function(t) {
          start.r = r(t);
          start.g = g(t);
          start.b = b(t);
          start.opacity = opacity(t);
          return start + "";
        };
      }

      rgb$1.gamma = rgbGamma;

      return rgb$1;
    })(1);

    function array(a, b) {
      var nb = b ? b.length : 0,
          na = a ? Math.min(nb, a.length) : 0,
          x = new Array(na),
          c = new Array(nb),
          i;

      for (i = 0; i < na; ++i) x[i] = interpolate(a[i], b[i]);
      for (; i < nb; ++i) c[i] = b[i];

      return function(t) {
        for (i = 0; i < na; ++i) c[i] = x[i](t);
        return c;
      };
    }

    function date(a, b) {
      var d = new Date;
      return a = +a, b -= a, function(t) {
        return d.setTime(a + b * t), d;
      };
    }

    function interpolateNumber(a, b) {
      return a = +a, b -= a, function(t) {
        return a + b * t;
      };
    }

    function object(a, b) {
      var i = {},
          c = {},
          k;

      if (a === null || typeof a !== "object") a = {};
      if (b === null || typeof b !== "object") b = {};

      for (k in b) {
        if (k in a) {
          i[k] = interpolate(a[k], b[k]);
        } else {
          c[k] = b[k];
        }
      }

      return function(t) {
        for (k in i) c[k] = i[k](t);
        return c;
      };
    }

    var reA = /[-+]?(?:\d+\.?\d*|\.?\d+)(?:[eE][-+]?\d+)?/g,
        reB = new RegExp(reA.source, "g");

    function zero(b) {
      return function() {
        return b;
      };
    }

    function one(b) {
      return function(t) {
        return b(t) + "";
      };
    }

    function string(a, b) {
      var bi = reA.lastIndex = reB.lastIndex = 0, // scan index for next number in b
          am, // current match in a
          bm, // current match in b
          bs, // string preceding current number in b, if any
          i = -1, // index in s
          s = [], // string constants and placeholders
          q = []; // number interpolators

      // Coerce inputs to strings.
      a = a + "", b = b + "";

      // Interpolate pairs of numbers in a & b.
      while ((am = reA.exec(a))
          && (bm = reB.exec(b))) {
        if ((bs = bm.index) > bi) { // a string precedes the next number in b
          bs = b.slice(bi, bs);
          if (s[i]) s[i] += bs; // coalesce with previous string
          else s[++i] = bs;
        }
        if ((am = am[0]) === (bm = bm[0])) { // numbers in a & b match
          if (s[i]) s[i] += bm; // coalesce with previous string
          else s[++i] = bm;
        } else { // interpolate non-matching numbers
          s[++i] = null;
          q.push({i: i, x: interpolateNumber(am, bm)});
        }
        bi = reB.lastIndex;
      }

      // Add remains of b.
      if (bi < b.length) {
        bs = b.slice(bi);
        if (s[i]) s[i] += bs; // coalesce with previous string
        else s[++i] = bs;
      }

      // Special optimization for only a single match.
      // Otherwise, interpolate each of the numbers and rejoin the string.
      return s.length < 2 ? (q[0]
          ? one(q[0].x)
          : zero(b))
          : (b = q.length, function(t) {
              for (var i = 0, o; i < b; ++i) s[(o = q[i]).i] = o.x(t);
              return s.join("");
            });
    }

    function interpolate(a, b) {
      var t = typeof b, c;
      return b == null || t === "boolean" ? constant(b)
          : (t === "number" ? interpolateNumber
          : t === "string" ? ((c = color(b)) ? (b = c, rgb$1) : string)
          : b instanceof color ? rgb$1
          : b instanceof Date ? date
          : Array.isArray(b) ? array
          : typeof b.valueOf !== "function" && typeof b.toString !== "function" || isNaN(b) ? object
          : interpolateNumber)(a, b);
    }

    function interpolateRound(a, b) {
      return a = +a, b -= a, function(t) {
        return Math.round(a + b * t);
      };
    }

    function constant$1(x) {
      return function() {
        return x;
      };
    }

    function number(x) {
      return +x;
    }

    var unit = [0, 1];

    function identity$1(x) {
      return x;
    }

    function normalize(a, b) {
      return (b -= (a = +a))
          ? function(x) { return (x - a) / b; }
          : constant$1(isNaN(b) ? NaN : 0.5);
    }

    function clamper(a, b) {
      var t;
      if (a > b) t = a, a = b, b = t;
      return function(x) { return Math.max(a, Math.min(b, x)); };
    }

    // normalize(a, b)(x) takes a domain value x in [a,b] and returns the corresponding parameter t in [0,1].
    // interpolate(a, b)(t) takes a parameter t in [0,1] and returns the corresponding range value x in [a,b].
    function bimap(domain, range, interpolate) {
      var d0 = domain[0], d1 = domain[1], r0 = range[0], r1 = range[1];
      if (d1 < d0) d0 = normalize(d1, d0), r0 = interpolate(r1, r0);
      else d0 = normalize(d0, d1), r0 = interpolate(r0, r1);
      return function(x) { return r0(d0(x)); };
    }

    function polymap(domain, range, interpolate) {
      var j = Math.min(domain.length, range.length) - 1,
          d = new Array(j),
          r = new Array(j),
          i = -1;

      // Reverse descending domains.
      if (domain[j] < domain[0]) {
        domain = domain.slice().reverse();
        range = range.slice().reverse();
      }

      while (++i < j) {
        d[i] = normalize(domain[i], domain[i + 1]);
        r[i] = interpolate(range[i], range[i + 1]);
      }

      return function(x) {
        var i = bisectRight(domain, x, 1, j) - 1;
        return r[i](d[i](x));
      };
    }

    function copy(source, target) {
      return target
          .domain(source.domain())
          .range(source.range())
          .interpolate(source.interpolate())
          .clamp(source.clamp())
          .unknown(source.unknown());
    }

    function transformer() {
      var domain = unit,
          range = unit,
          interpolate$1 = interpolate,
          transform,
          untransform,
          unknown,
          clamp = identity$1,
          piecewise,
          output,
          input;

      function rescale() {
        var n = Math.min(domain.length, range.length);
        if (clamp !== identity$1) clamp = clamper(domain[0], domain[n - 1]);
        piecewise = n > 2 ? polymap : bimap;
        output = input = null;
        return scale;
      }

      function scale(x) {
        return isNaN(x = +x) ? unknown : (output || (output = piecewise(domain.map(transform), range, interpolate$1)))(transform(clamp(x)));
      }

      scale.invert = function(y) {
        return clamp(untransform((input || (input = piecewise(range, domain.map(transform), interpolateNumber)))(y)));
      };

      scale.domain = function(_) {
        return arguments.length ? (domain = Array.from(_, number), rescale()) : domain.slice();
      };

      scale.range = function(_) {
        return arguments.length ? (range = Array.from(_), rescale()) : range.slice();
      };

      scale.rangeRound = function(_) {
        return range = Array.from(_), interpolate$1 = interpolateRound, rescale();
      };

      scale.clamp = function(_) {
        return arguments.length ? (clamp = _ ? true : identity$1, rescale()) : clamp !== identity$1;
      };

      scale.interpolate = function(_) {
        return arguments.length ? (interpolate$1 = _, rescale()) : interpolate$1;
      };

      scale.unknown = function(_) {
        return arguments.length ? (unknown = _, scale) : unknown;
      };

      return function(t, u) {
        transform = t, untransform = u;
        return rescale();
      };
    }

    function continuous() {
      return transformer()(identity$1, identity$1);
    }

    // Computes the decimal coefficient and exponent of the specified number x with
    // significant digits p, where x is positive and p is in [1, 21] or undefined.
    // For example, formatDecimal(1.23) returns ["123", 0].
    function formatDecimal(x, p) {
      if ((i = (x = p ? x.toExponential(p - 1) : x.toExponential()).indexOf("e")) < 0) return null; // NaN, ±Infinity
      var i, coefficient = x.slice(0, i);

      // The string returned by toExponential either has the form \d\.\d+e[-+]\d+
      // (e.g., 1.2e+3) or the form \de[-+]\d+ (e.g., 1e+3).
      return [
        coefficient.length > 1 ? coefficient[0] + coefficient.slice(2) : coefficient,
        +x.slice(i + 1)
      ];
    }

    function exponent(x) {
      return x = formatDecimal(Math.abs(x)), x ? x[1] : NaN;
    }

    function formatGroup(grouping, thousands) {
      return function(value, width) {
        var i = value.length,
            t = [],
            j = 0,
            g = grouping[0],
            length = 0;

        while (i > 0 && g > 0) {
          if (length + g + 1 > width) g = Math.max(1, width - length);
          t.push(value.substring(i -= g, i + g));
          if ((length += g + 1) > width) break;
          g = grouping[j = (j + 1) % grouping.length];
        }

        return t.reverse().join(thousands);
      };
    }

    function formatNumerals(numerals) {
      return function(value) {
        return value.replace(/[0-9]/g, function(i) {
          return numerals[+i];
        });
      };
    }

    // [[fill]align][sign][symbol][0][width][,][.precision][~][type]
    var re = /^(?:(.)?([<>=^]))?([+\-( ])?([$#])?(0)?(\d+)?(,)?(\.\d+)?(~)?([a-z%])?$/i;

    function formatSpecifier(specifier) {
      if (!(match = re.exec(specifier))) throw new Error("invalid format: " + specifier);
      var match;
      return new FormatSpecifier({
        fill: match[1],
        align: match[2],
        sign: match[3],
        symbol: match[4],
        zero: match[5],
        width: match[6],
        comma: match[7],
        precision: match[8] && match[8].slice(1),
        trim: match[9],
        type: match[10]
      });
    }

    formatSpecifier.prototype = FormatSpecifier.prototype; // instanceof

    function FormatSpecifier(specifier) {
      this.fill = specifier.fill === undefined ? " " : specifier.fill + "";
      this.align = specifier.align === undefined ? ">" : specifier.align + "";
      this.sign = specifier.sign === undefined ? "-" : specifier.sign + "";
      this.symbol = specifier.symbol === undefined ? "" : specifier.symbol + "";
      this.zero = !!specifier.zero;
      this.width = specifier.width === undefined ? undefined : +specifier.width;
      this.comma = !!specifier.comma;
      this.precision = specifier.precision === undefined ? undefined : +specifier.precision;
      this.trim = !!specifier.trim;
      this.type = specifier.type === undefined ? "" : specifier.type + "";
    }

    FormatSpecifier.prototype.toString = function() {
      return this.fill
          + this.align
          + this.sign
          + this.symbol
          + (this.zero ? "0" : "")
          + (this.width === undefined ? "" : Math.max(1, this.width | 0))
          + (this.comma ? "," : "")
          + (this.precision === undefined ? "" : "." + Math.max(0, this.precision | 0))
          + (this.trim ? "~" : "")
          + this.type;
    };

    // Trims insignificant zeros, e.g., replaces 1.2000k with 1.2k.
    function formatTrim(s) {
      out: for (var n = s.length, i = 1, i0 = -1, i1; i < n; ++i) {
        switch (s[i]) {
          case ".": i0 = i1 = i; break;
          case "0": if (i0 === 0) i0 = i; i1 = i; break;
          default: if (i0 > 0) { if (!+s[i]) break out; i0 = 0; } break;
        }
      }
      return i0 > 0 ? s.slice(0, i0) + s.slice(i1 + 1) : s;
    }

    var prefixExponent;

    function formatPrefixAuto(x, p) {
      var d = formatDecimal(x, p);
      if (!d) return x + "";
      var coefficient = d[0],
          exponent = d[1],
          i = exponent - (prefixExponent = Math.max(-8, Math.min(8, Math.floor(exponent / 3))) * 3) + 1,
          n = coefficient.length;
      return i === n ? coefficient
          : i > n ? coefficient + new Array(i - n + 1).join("0")
          : i > 0 ? coefficient.slice(0, i) + "." + coefficient.slice(i)
          : "0." + new Array(1 - i).join("0") + formatDecimal(x, Math.max(0, p + i - 1))[0]; // less than 1y!
    }

    function formatRounded(x, p) {
      var d = formatDecimal(x, p);
      if (!d) return x + "";
      var coefficient = d[0],
          exponent = d[1];
      return exponent < 0 ? "0." + new Array(-exponent).join("0") + coefficient
          : coefficient.length > exponent + 1 ? coefficient.slice(0, exponent + 1) + "." + coefficient.slice(exponent + 1)
          : coefficient + new Array(exponent - coefficient.length + 2).join("0");
    }

    var formatTypes = {
      "%": function(x, p) { return (x * 100).toFixed(p); },
      "b": function(x) { return Math.round(x).toString(2); },
      "c": function(x) { return x + ""; },
      "d": function(x) { return Math.round(x).toString(10); },
      "e": function(x, p) { return x.toExponential(p); },
      "f": function(x, p) { return x.toFixed(p); },
      "g": function(x, p) { return x.toPrecision(p); },
      "o": function(x) { return Math.round(x).toString(8); },
      "p": function(x, p) { return formatRounded(x * 100, p); },
      "r": formatRounded,
      "s": formatPrefixAuto,
      "X": function(x) { return Math.round(x).toString(16).toUpperCase(); },
      "x": function(x) { return Math.round(x).toString(16); }
    };

    function identity$2(x) {
      return x;
    }

    var map = Array.prototype.map,
        prefixes = ["y","z","a","f","p","n","µ","m","","k","M","G","T","P","E","Z","Y"];

    function formatLocale(locale) {
      var group = locale.grouping === undefined || locale.thousands === undefined ? identity$2 : formatGroup(map.call(locale.grouping, Number), locale.thousands + ""),
          currencyPrefix = locale.currency === undefined ? "" : locale.currency[0] + "",
          currencySuffix = locale.currency === undefined ? "" : locale.currency[1] + "",
          decimal = locale.decimal === undefined ? "." : locale.decimal + "",
          numerals = locale.numerals === undefined ? identity$2 : formatNumerals(map.call(locale.numerals, String)),
          percent = locale.percent === undefined ? "%" : locale.percent + "",
          minus = locale.minus === undefined ? "-" : locale.minus + "",
          nan = locale.nan === undefined ? "NaN" : locale.nan + "";

      function newFormat(specifier) {
        specifier = formatSpecifier(specifier);

        var fill = specifier.fill,
            align = specifier.align,
            sign = specifier.sign,
            symbol = specifier.symbol,
            zero = specifier.zero,
            width = specifier.width,
            comma = specifier.comma,
            precision = specifier.precision,
            trim = specifier.trim,
            type = specifier.type;

        // The "n" type is an alias for ",g".
        if (type === "n") comma = true, type = "g";

        // The "" type, and any invalid type, is an alias for ".12~g".
        else if (!formatTypes[type]) precision === undefined && (precision = 12), trim = true, type = "g";

        // If zero fill is specified, padding goes after sign and before digits.
        if (zero || (fill === "0" && align === "=")) zero = true, fill = "0", align = "=";

        // Compute the prefix and suffix.
        // For SI-prefix, the suffix is lazily computed.
        var prefix = symbol === "$" ? currencyPrefix : symbol === "#" && /[boxX]/.test(type) ? "0" + type.toLowerCase() : "",
            suffix = symbol === "$" ? currencySuffix : /[%p]/.test(type) ? percent : "";

        // What format function should we use?
        // Is this an integer type?
        // Can this type generate exponential notation?
        var formatType = formatTypes[type],
            maybeSuffix = /[defgprs%]/.test(type);

        // Set the default precision if not specified,
        // or clamp the specified precision to the supported range.
        // For significant precision, it must be in [1, 21].
        // For fixed precision, it must be in [0, 20].
        precision = precision === undefined ? 6
            : /[gprs]/.test(type) ? Math.max(1, Math.min(21, precision))
            : Math.max(0, Math.min(20, precision));

        function format(value) {
          var valuePrefix = prefix,
              valueSuffix = suffix,
              i, n, c;

          if (type === "c") {
            valueSuffix = formatType(value) + valueSuffix;
            value = "";
          } else {
            value = +value;

            // Perform the initial formatting.
            var valueNegative = value < 0;
            value = isNaN(value) ? nan : formatType(Math.abs(value), precision);

            // Trim insignificant zeros.
            if (trim) value = formatTrim(value);

            // If a negative value rounds to zero during formatting, treat as positive.
            if (valueNegative && +value === 0) valueNegative = false;

            // Compute the prefix and suffix.
            valuePrefix = (valueNegative ? (sign === "(" ? sign : minus) : sign === "-" || sign === "(" ? "" : sign) + valuePrefix;

            valueSuffix = (type === "s" ? prefixes[8 + prefixExponent / 3] : "") + valueSuffix + (valueNegative && sign === "(" ? ")" : "");

            // Break the formatted value into the integer “value” part that can be
            // grouped, and fractional or exponential “suffix” part that is not.
            if (maybeSuffix) {
              i = -1, n = value.length;
              while (++i < n) {
                if (c = value.charCodeAt(i), 48 > c || c > 57) {
                  valueSuffix = (c === 46 ? decimal + value.slice(i + 1) : value.slice(i)) + valueSuffix;
                  value = value.slice(0, i);
                  break;
                }
              }
            }
          }

          // If the fill character is not "0", grouping is applied before padding.
          if (comma && !zero) value = group(value, Infinity);

          // Compute the padding.
          var length = valuePrefix.length + value.length + valueSuffix.length,
              padding = length < width ? new Array(width - length + 1).join(fill) : "";

          // If the fill character is "0", grouping is applied after padding.
          if (comma && zero) value = group(padding + value, padding.length ? width - valueSuffix.length : Infinity), padding = "";

          // Reconstruct the final output based on the desired alignment.
          switch (align) {
            case "<": value = valuePrefix + value + valueSuffix + padding; break;
            case "=": value = valuePrefix + padding + value + valueSuffix; break;
            case "^": value = padding.slice(0, length = padding.length >> 1) + valuePrefix + value + valueSuffix + padding.slice(length); break;
            default: value = padding + valuePrefix + value + valueSuffix; break;
          }

          return numerals(value);
        }

        format.toString = function() {
          return specifier + "";
        };

        return format;
      }

      function formatPrefix(specifier, value) {
        var f = newFormat((specifier = formatSpecifier(specifier), specifier.type = "f", specifier)),
            e = Math.max(-8, Math.min(8, Math.floor(exponent(value) / 3))) * 3,
            k = Math.pow(10, -e),
            prefix = prefixes[8 + e / 3];
        return function(value) {
          return f(k * value) + prefix;
        };
      }

      return {
        format: newFormat,
        formatPrefix: formatPrefix
      };
    }

    var locale;
    var format;
    var formatPrefix;

    defaultLocale({
      decimal: ".",
      thousands: ",",
      grouping: [3],
      currency: ["$", ""],
      minus: "-"
    });

    function defaultLocale(definition) {
      locale = formatLocale(definition);
      format = locale.format;
      formatPrefix = locale.formatPrefix;
      return locale;
    }

    function precisionFixed(step) {
      return Math.max(0, -exponent(Math.abs(step)));
    }

    function precisionPrefix(step, value) {
      return Math.max(0, Math.max(-8, Math.min(8, Math.floor(exponent(value) / 3))) * 3 - exponent(Math.abs(step)));
    }

    function precisionRound(step, max) {
      step = Math.abs(step), max = Math.abs(max) - step;
      return Math.max(0, exponent(max) - exponent(step)) + 1;
    }

    function tickFormat(start, stop, count, specifier) {
      var step = tickStep(start, stop, count),
          precision;
      specifier = formatSpecifier(specifier == null ? ",f" : specifier);
      switch (specifier.type) {
        case "s": {
          var value = Math.max(Math.abs(start), Math.abs(stop));
          if (specifier.precision == null && !isNaN(precision = precisionPrefix(step, value))) specifier.precision = precision;
          return formatPrefix(specifier, value);
        }
        case "":
        case "e":
        case "g":
        case "p":
        case "r": {
          if (specifier.precision == null && !isNaN(precision = precisionRound(step, Math.max(Math.abs(start), Math.abs(stop))))) specifier.precision = precision - (specifier.type === "e");
          break;
        }
        case "f":
        case "%": {
          if (specifier.precision == null && !isNaN(precision = precisionFixed(step))) specifier.precision = precision - (specifier.type === "%") * 2;
          break;
        }
      }
      return format(specifier);
    }

    function linearish(scale) {
      var domain = scale.domain;

      scale.ticks = function(count) {
        var d = domain();
        return ticks(d[0], d[d.length - 1], count == null ? 10 : count);
      };

      scale.tickFormat = function(count, specifier) {
        var d = domain();
        return tickFormat(d[0], d[d.length - 1], count == null ? 10 : count, specifier);
      };

      scale.nice = function(count) {
        if (count == null) count = 10;

        var d = domain(),
            i0 = 0,
            i1 = d.length - 1,
            start = d[i0],
            stop = d[i1],
            step;

        if (stop < start) {
          step = start, start = stop, stop = step;
          step = i0, i0 = i1, i1 = step;
        }

        step = tickIncrement(start, stop, count);

        if (step > 0) {
          start = Math.floor(start / step) * step;
          stop = Math.ceil(stop / step) * step;
          step = tickIncrement(start, stop, count);
        } else if (step < 0) {
          start = Math.ceil(start * step) / step;
          stop = Math.floor(stop * step) / step;
          step = tickIncrement(start, stop, count);
        }

        if (step > 0) {
          d[i0] = Math.floor(start / step) * step;
          d[i1] = Math.ceil(stop / step) * step;
          domain(d);
        } else if (step < 0) {
          d[i0] = Math.ceil(start * step) / step;
          d[i1] = Math.floor(stop * step) / step;
          domain(d);
        }

        return scale;
      };

      return scale;
    }

    function linear$1() {
      var scale = continuous();

      scale.copy = function() {
        return copy(scale, linear$1());
      };

      initRange.apply(scale, arguments);

      return linearish(scale);
    }

    var t0$1 = new Date,
        t1$1 = new Date;

    function newInterval(floori, offseti, count, field) {

      function interval(date) {
        return floori(date = arguments.length === 0 ? new Date : new Date(+date)), date;
      }

      interval.floor = function(date) {
        return floori(date = new Date(+date)), date;
      };

      interval.ceil = function(date) {
        return floori(date = new Date(date - 1)), offseti(date, 1), floori(date), date;
      };

      interval.round = function(date) {
        var d0 = interval(date),
            d1 = interval.ceil(date);
        return date - d0 < d1 - date ? d0 : d1;
      };

      interval.offset = function(date, step) {
        return offseti(date = new Date(+date), step == null ? 1 : Math.floor(step)), date;
      };

      interval.range = function(start, stop, step) {
        var range = [], previous;
        start = interval.ceil(start);
        step = step == null ? 1 : Math.floor(step);
        if (!(start < stop) || !(step > 0)) return range; // also handles Invalid Date
        do range.push(previous = new Date(+start)), offseti(start, step), floori(start);
        while (previous < start && start < stop);
        return range;
      };

      interval.filter = function(test) {
        return newInterval(function(date) {
          if (date >= date) while (floori(date), !test(date)) date.setTime(date - 1);
        }, function(date, step) {
          if (date >= date) {
            if (step < 0) while (++step <= 0) {
              while (offseti(date, -1), !test(date)) {} // eslint-disable-line no-empty
            } else while (--step >= 0) {
              while (offseti(date, +1), !test(date)) {} // eslint-disable-line no-empty
            }
          }
        });
      };

      if (count) {
        interval.count = function(start, end) {
          t0$1.setTime(+start), t1$1.setTime(+end);
          floori(t0$1), floori(t1$1);
          return Math.floor(count(t0$1, t1$1));
        };

        interval.every = function(step) {
          step = Math.floor(step);
          return !isFinite(step) || !(step > 0) ? null
              : !(step > 1) ? interval
              : interval.filter(field
                  ? function(d) { return field(d) % step === 0; }
                  : function(d) { return interval.count(0, d) % step === 0; });
        };
      }

      return interval;
    }

    var millisecond = newInterval(function() {
      // noop
    }, function(date, step) {
      date.setTime(+date + step);
    }, function(start, end) {
      return end - start;
    });

    // An optimized implementation for this simple case.
    millisecond.every = function(k) {
      k = Math.floor(k);
      if (!isFinite(k) || !(k > 0)) return null;
      if (!(k > 1)) return millisecond;
      return newInterval(function(date) {
        date.setTime(Math.floor(date / k) * k);
      }, function(date, step) {
        date.setTime(+date + step * k);
      }, function(start, end) {
        return (end - start) / k;
      });
    };

    var durationSecond = 1e3;
    var durationMinute = 6e4;
    var durationHour = 36e5;
    var durationDay = 864e5;
    var durationWeek = 6048e5;

    var second = newInterval(function(date) {
      date.setTime(date - date.getMilliseconds());
    }, function(date, step) {
      date.setTime(+date + step * durationSecond);
    }, function(start, end) {
      return (end - start) / durationSecond;
    }, function(date) {
      return date.getUTCSeconds();
    });

    var minute = newInterval(function(date) {
      date.setTime(date - date.getMilliseconds() - date.getSeconds() * durationSecond);
    }, function(date, step) {
      date.setTime(+date + step * durationMinute);
    }, function(start, end) {
      return (end - start) / durationMinute;
    }, function(date) {
      return date.getMinutes();
    });

    var hour = newInterval(function(date) {
      date.setTime(date - date.getMilliseconds() - date.getSeconds() * durationSecond - date.getMinutes() * durationMinute);
    }, function(date, step) {
      date.setTime(+date + step * durationHour);
    }, function(start, end) {
      return (end - start) / durationHour;
    }, function(date) {
      return date.getHours();
    });

    var day = newInterval(function(date) {
      date.setHours(0, 0, 0, 0);
    }, function(date, step) {
      date.setDate(date.getDate() + step);
    }, function(start, end) {
      return (end - start - (end.getTimezoneOffset() - start.getTimezoneOffset()) * durationMinute) / durationDay;
    }, function(date) {
      return date.getDate() - 1;
    });

    function weekday(i) {
      return newInterval(function(date) {
        date.setDate(date.getDate() - (date.getDay() + 7 - i) % 7);
        date.setHours(0, 0, 0, 0);
      }, function(date, step) {
        date.setDate(date.getDate() + step * 7);
      }, function(start, end) {
        return (end - start - (end.getTimezoneOffset() - start.getTimezoneOffset()) * durationMinute) / durationWeek;
      });
    }

    var sunday = weekday(0);
    var monday = weekday(1);
    var tuesday = weekday(2);
    var wednesday = weekday(3);
    var thursday = weekday(4);
    var friday = weekday(5);
    var saturday = weekday(6);

    var month = newInterval(function(date) {
      date.setDate(1);
      date.setHours(0, 0, 0, 0);
    }, function(date, step) {
      date.setMonth(date.getMonth() + step);
    }, function(start, end) {
      return end.getMonth() - start.getMonth() + (end.getFullYear() - start.getFullYear()) * 12;
    }, function(date) {
      return date.getMonth();
    });

    var year = newInterval(function(date) {
      date.setMonth(0, 1);
      date.setHours(0, 0, 0, 0);
    }, function(date, step) {
      date.setFullYear(date.getFullYear() + step);
    }, function(start, end) {
      return end.getFullYear() - start.getFullYear();
    }, function(date) {
      return date.getFullYear();
    });

    // An optimized implementation for this simple case.
    year.every = function(k) {
      return !isFinite(k = Math.floor(k)) || !(k > 0) ? null : newInterval(function(date) {
        date.setFullYear(Math.floor(date.getFullYear() / k) * k);
        date.setMonth(0, 1);
        date.setHours(0, 0, 0, 0);
      }, function(date, step) {
        date.setFullYear(date.getFullYear() + step * k);
      });
    };

    var utcMinute = newInterval(function(date) {
      date.setUTCSeconds(0, 0);
    }, function(date, step) {
      date.setTime(+date + step * durationMinute);
    }, function(start, end) {
      return (end - start) / durationMinute;
    }, function(date) {
      return date.getUTCMinutes();
    });

    var utcHour = newInterval(function(date) {
      date.setUTCMinutes(0, 0, 0);
    }, function(date, step) {
      date.setTime(+date + step * durationHour);
    }, function(start, end) {
      return (end - start) / durationHour;
    }, function(date) {
      return date.getUTCHours();
    });

    var utcDay = newInterval(function(date) {
      date.setUTCHours(0, 0, 0, 0);
    }, function(date, step) {
      date.setUTCDate(date.getUTCDate() + step);
    }, function(start, end) {
      return (end - start) / durationDay;
    }, function(date) {
      return date.getUTCDate() - 1;
    });

    function utcWeekday(i) {
      return newInterval(function(date) {
        date.setUTCDate(date.getUTCDate() - (date.getUTCDay() + 7 - i) % 7);
        date.setUTCHours(0, 0, 0, 0);
      }, function(date, step) {
        date.setUTCDate(date.getUTCDate() + step * 7);
      }, function(start, end) {
        return (end - start) / durationWeek;
      });
    }

    var utcSunday = utcWeekday(0);
    var utcMonday = utcWeekday(1);
    var utcTuesday = utcWeekday(2);
    var utcWednesday = utcWeekday(3);
    var utcThursday = utcWeekday(4);
    var utcFriday = utcWeekday(5);
    var utcSaturday = utcWeekday(6);

    var utcMonth = newInterval(function(date) {
      date.setUTCDate(1);
      date.setUTCHours(0, 0, 0, 0);
    }, function(date, step) {
      date.setUTCMonth(date.getUTCMonth() + step);
    }, function(start, end) {
      return end.getUTCMonth() - start.getUTCMonth() + (end.getUTCFullYear() - start.getUTCFullYear()) * 12;
    }, function(date) {
      return date.getUTCMonth();
    });

    var utcYear = newInterval(function(date) {
      date.setUTCMonth(0, 1);
      date.setUTCHours(0, 0, 0, 0);
    }, function(date, step) {
      date.setUTCFullYear(date.getUTCFullYear() + step);
    }, function(start, end) {
      return end.getUTCFullYear() - start.getUTCFullYear();
    }, function(date) {
      return date.getUTCFullYear();
    });

    // An optimized implementation for this simple case.
    utcYear.every = function(k) {
      return !isFinite(k = Math.floor(k)) || !(k > 0) ? null : newInterval(function(date) {
        date.setUTCFullYear(Math.floor(date.getUTCFullYear() / k) * k);
        date.setUTCMonth(0, 1);
        date.setUTCHours(0, 0, 0, 0);
      }, function(date, step) {
        date.setUTCFullYear(date.getUTCFullYear() + step * k);
      });
    };

    function localDate(d) {
      if (0 <= d.y && d.y < 100) {
        var date = new Date(-1, d.m, d.d, d.H, d.M, d.S, d.L);
        date.setFullYear(d.y);
        return date;
      }
      return new Date(d.y, d.m, d.d, d.H, d.M, d.S, d.L);
    }

    function utcDate(d) {
      if (0 <= d.y && d.y < 100) {
        var date = new Date(Date.UTC(-1, d.m, d.d, d.H, d.M, d.S, d.L));
        date.setUTCFullYear(d.y);
        return date;
      }
      return new Date(Date.UTC(d.y, d.m, d.d, d.H, d.M, d.S, d.L));
    }

    function newDate(y, m, d) {
      return {y: y, m: m, d: d, H: 0, M: 0, S: 0, L: 0};
    }

    function formatLocale$1(locale) {
      var locale_dateTime = locale.dateTime,
          locale_date = locale.date,
          locale_time = locale.time,
          locale_periods = locale.periods,
          locale_weekdays = locale.days,
          locale_shortWeekdays = locale.shortDays,
          locale_months = locale.months,
          locale_shortMonths = locale.shortMonths;

      var periodRe = formatRe(locale_periods),
          periodLookup = formatLookup(locale_periods),
          weekdayRe = formatRe(locale_weekdays),
          weekdayLookup = formatLookup(locale_weekdays),
          shortWeekdayRe = formatRe(locale_shortWeekdays),
          shortWeekdayLookup = formatLookup(locale_shortWeekdays),
          monthRe = formatRe(locale_months),
          monthLookup = formatLookup(locale_months),
          shortMonthRe = formatRe(locale_shortMonths),
          shortMonthLookup = formatLookup(locale_shortMonths);

      var formats = {
        "a": formatShortWeekday,
        "A": formatWeekday,
        "b": formatShortMonth,
        "B": formatMonth,
        "c": null,
        "d": formatDayOfMonth,
        "e": formatDayOfMonth,
        "f": formatMicroseconds,
        "H": formatHour24,
        "I": formatHour12,
        "j": formatDayOfYear,
        "L": formatMilliseconds,
        "m": formatMonthNumber,
        "M": formatMinutes,
        "p": formatPeriod,
        "q": formatQuarter,
        "Q": formatUnixTimestamp,
        "s": formatUnixTimestampSeconds,
        "S": formatSeconds,
        "u": formatWeekdayNumberMonday,
        "U": formatWeekNumberSunday,
        "V": formatWeekNumberISO,
        "w": formatWeekdayNumberSunday,
        "W": formatWeekNumberMonday,
        "x": null,
        "X": null,
        "y": formatYear,
        "Y": formatFullYear,
        "Z": formatZone,
        "%": formatLiteralPercent
      };

      var utcFormats = {
        "a": formatUTCShortWeekday,
        "A": formatUTCWeekday,
        "b": formatUTCShortMonth,
        "B": formatUTCMonth,
        "c": null,
        "d": formatUTCDayOfMonth,
        "e": formatUTCDayOfMonth,
        "f": formatUTCMicroseconds,
        "H": formatUTCHour24,
        "I": formatUTCHour12,
        "j": formatUTCDayOfYear,
        "L": formatUTCMilliseconds,
        "m": formatUTCMonthNumber,
        "M": formatUTCMinutes,
        "p": formatUTCPeriod,
        "q": formatUTCQuarter,
        "Q": formatUnixTimestamp,
        "s": formatUnixTimestampSeconds,
        "S": formatUTCSeconds,
        "u": formatUTCWeekdayNumberMonday,
        "U": formatUTCWeekNumberSunday,
        "V": formatUTCWeekNumberISO,
        "w": formatUTCWeekdayNumberSunday,
        "W": formatUTCWeekNumberMonday,
        "x": null,
        "X": null,
        "y": formatUTCYear,
        "Y": formatUTCFullYear,
        "Z": formatUTCZone,
        "%": formatLiteralPercent
      };

      var parses = {
        "a": parseShortWeekday,
        "A": parseWeekday,
        "b": parseShortMonth,
        "B": parseMonth,
        "c": parseLocaleDateTime,
        "d": parseDayOfMonth,
        "e": parseDayOfMonth,
        "f": parseMicroseconds,
        "H": parseHour24,
        "I": parseHour24,
        "j": parseDayOfYear,
        "L": parseMilliseconds,
        "m": parseMonthNumber,
        "M": parseMinutes,
        "p": parsePeriod,
        "q": parseQuarter,
        "Q": parseUnixTimestamp,
        "s": parseUnixTimestampSeconds,
        "S": parseSeconds,
        "u": parseWeekdayNumberMonday,
        "U": parseWeekNumberSunday,
        "V": parseWeekNumberISO,
        "w": parseWeekdayNumberSunday,
        "W": parseWeekNumberMonday,
        "x": parseLocaleDate,
        "X": parseLocaleTime,
        "y": parseYear,
        "Y": parseFullYear,
        "Z": parseZone,
        "%": parseLiteralPercent
      };

      // These recursive directive definitions must be deferred.
      formats.x = newFormat(locale_date, formats);
      formats.X = newFormat(locale_time, formats);
      formats.c = newFormat(locale_dateTime, formats);
      utcFormats.x = newFormat(locale_date, utcFormats);
      utcFormats.X = newFormat(locale_time, utcFormats);
      utcFormats.c = newFormat(locale_dateTime, utcFormats);

      function newFormat(specifier, formats) {
        return function(date) {
          var string = [],
              i = -1,
              j = 0,
              n = specifier.length,
              c,
              pad,
              format;

          if (!(date instanceof Date)) date = new Date(+date);

          while (++i < n) {
            if (specifier.charCodeAt(i) === 37) {
              string.push(specifier.slice(j, i));
              if ((pad = pads[c = specifier.charAt(++i)]) != null) c = specifier.charAt(++i);
              else pad = c === "e" ? " " : "0";
              if (format = formats[c]) c = format(date, pad);
              string.push(c);
              j = i + 1;
            }
          }

          string.push(specifier.slice(j, i));
          return string.join("");
        };
      }

      function newParse(specifier, Z) {
        return function(string) {
          var d = newDate(1900, undefined, 1),
              i = parseSpecifier(d, specifier, string += "", 0),
              week, day$1;
          if (i != string.length) return null;

          // If a UNIX timestamp is specified, return it.
          if ("Q" in d) return new Date(d.Q);
          if ("s" in d) return new Date(d.s * 1000 + ("L" in d ? d.L : 0));

          // If this is utcParse, never use the local timezone.
          if (Z && !("Z" in d)) d.Z = 0;

          // The am-pm flag is 0 for AM, and 1 for PM.
          if ("p" in d) d.H = d.H % 12 + d.p * 12;

          // If the month was not specified, inherit from the quarter.
          if (d.m === undefined) d.m = "q" in d ? d.q : 0;

          // Convert day-of-week and week-of-year to day-of-year.
          if ("V" in d) {
            if (d.V < 1 || d.V > 53) return null;
            if (!("w" in d)) d.w = 1;
            if ("Z" in d) {
              week = utcDate(newDate(d.y, 0, 1)), day$1 = week.getUTCDay();
              week = day$1 > 4 || day$1 === 0 ? utcMonday.ceil(week) : utcMonday(week);
              week = utcDay.offset(week, (d.V - 1) * 7);
              d.y = week.getUTCFullYear();
              d.m = week.getUTCMonth();
              d.d = week.getUTCDate() + (d.w + 6) % 7;
            } else {
              week = localDate(newDate(d.y, 0, 1)), day$1 = week.getDay();
              week = day$1 > 4 || day$1 === 0 ? monday.ceil(week) : monday(week);
              week = day.offset(week, (d.V - 1) * 7);
              d.y = week.getFullYear();
              d.m = week.getMonth();
              d.d = week.getDate() + (d.w + 6) % 7;
            }
          } else if ("W" in d || "U" in d) {
            if (!("w" in d)) d.w = "u" in d ? d.u % 7 : "W" in d ? 1 : 0;
            day$1 = "Z" in d ? utcDate(newDate(d.y, 0, 1)).getUTCDay() : localDate(newDate(d.y, 0, 1)).getDay();
            d.m = 0;
            d.d = "W" in d ? (d.w + 6) % 7 + d.W * 7 - (day$1 + 5) % 7 : d.w + d.U * 7 - (day$1 + 6) % 7;
          }

          // If a time zone is specified, all fields are interpreted as UTC and then
          // offset according to the specified time zone.
          if ("Z" in d) {
            d.H += d.Z / 100 | 0;
            d.M += d.Z % 100;
            return utcDate(d);
          }

          // Otherwise, all fields are in local time.
          return localDate(d);
        };
      }

      function parseSpecifier(d, specifier, string, j) {
        var i = 0,
            n = specifier.length,
            m = string.length,
            c,
            parse;

        while (i < n) {
          if (j >= m) return -1;
          c = specifier.charCodeAt(i++);
          if (c === 37) {
            c = specifier.charAt(i++);
            parse = parses[c in pads ? specifier.charAt(i++) : c];
            if (!parse || ((j = parse(d, string, j)) < 0)) return -1;
          } else if (c != string.charCodeAt(j++)) {
            return -1;
          }
        }

        return j;
      }

      function parsePeriod(d, string, i) {
        var n = periodRe.exec(string.slice(i));
        return n ? (d.p = periodLookup[n[0].toLowerCase()], i + n[0].length) : -1;
      }

      function parseShortWeekday(d, string, i) {
        var n = shortWeekdayRe.exec(string.slice(i));
        return n ? (d.w = shortWeekdayLookup[n[0].toLowerCase()], i + n[0].length) : -1;
      }

      function parseWeekday(d, string, i) {
        var n = weekdayRe.exec(string.slice(i));
        return n ? (d.w = weekdayLookup[n[0].toLowerCase()], i + n[0].length) : -1;
      }

      function parseShortMonth(d, string, i) {
        var n = shortMonthRe.exec(string.slice(i));
        return n ? (d.m = shortMonthLookup[n[0].toLowerCase()], i + n[0].length) : -1;
      }

      function parseMonth(d, string, i) {
        var n = monthRe.exec(string.slice(i));
        return n ? (d.m = monthLookup[n[0].toLowerCase()], i + n[0].length) : -1;
      }

      function parseLocaleDateTime(d, string, i) {
        return parseSpecifier(d, locale_dateTime, string, i);
      }

      function parseLocaleDate(d, string, i) {
        return parseSpecifier(d, locale_date, string, i);
      }

      function parseLocaleTime(d, string, i) {
        return parseSpecifier(d, locale_time, string, i);
      }

      function formatShortWeekday(d) {
        return locale_shortWeekdays[d.getDay()];
      }

      function formatWeekday(d) {
        return locale_weekdays[d.getDay()];
      }

      function formatShortMonth(d) {
        return locale_shortMonths[d.getMonth()];
      }

      function formatMonth(d) {
        return locale_months[d.getMonth()];
      }

      function formatPeriod(d) {
        return locale_periods[+(d.getHours() >= 12)];
      }

      function formatQuarter(d) {
        return 1 + ~~(d.getMonth() / 3);
      }

      function formatUTCShortWeekday(d) {
        return locale_shortWeekdays[d.getUTCDay()];
      }

      function formatUTCWeekday(d) {
        return locale_weekdays[d.getUTCDay()];
      }

      function formatUTCShortMonth(d) {
        return locale_shortMonths[d.getUTCMonth()];
      }

      function formatUTCMonth(d) {
        return locale_months[d.getUTCMonth()];
      }

      function formatUTCPeriod(d) {
        return locale_periods[+(d.getUTCHours() >= 12)];
      }

      function formatUTCQuarter(d) {
        return 1 + ~~(d.getUTCMonth() / 3);
      }

      return {
        format: function(specifier) {
          var f = newFormat(specifier += "", formats);
          f.toString = function() { return specifier; };
          return f;
        },
        parse: function(specifier) {
          var p = newParse(specifier += "", false);
          p.toString = function() { return specifier; };
          return p;
        },
        utcFormat: function(specifier) {
          var f = newFormat(specifier += "", utcFormats);
          f.toString = function() { return specifier; };
          return f;
        },
        utcParse: function(specifier) {
          var p = newParse(specifier += "", true);
          p.toString = function() { return specifier; };
          return p;
        }
      };
    }

    var pads = {"-": "", "_": " ", "0": "0"},
        numberRe = /^\s*\d+/, // note: ignores next directive
        percentRe = /^%/,
        requoteRe = /[\\^$*+?|[\]().{}]/g;

    function pad(value, fill, width) {
      var sign = value < 0 ? "-" : "",
          string = (sign ? -value : value) + "",
          length = string.length;
      return sign + (length < width ? new Array(width - length + 1).join(fill) + string : string);
    }

    function requote(s) {
      return s.replace(requoteRe, "\\$&");
    }

    function formatRe(names) {
      return new RegExp("^(?:" + names.map(requote).join("|") + ")", "i");
    }

    function formatLookup(names) {
      var map = {}, i = -1, n = names.length;
      while (++i < n) map[names[i].toLowerCase()] = i;
      return map;
    }

    function parseWeekdayNumberSunday(d, string, i) {
      var n = numberRe.exec(string.slice(i, i + 1));
      return n ? (d.w = +n[0], i + n[0].length) : -1;
    }

    function parseWeekdayNumberMonday(d, string, i) {
      var n = numberRe.exec(string.slice(i, i + 1));
      return n ? (d.u = +n[0], i + n[0].length) : -1;
    }

    function parseWeekNumberSunday(d, string, i) {
      var n = numberRe.exec(string.slice(i, i + 2));
      return n ? (d.U = +n[0], i + n[0].length) : -1;
    }

    function parseWeekNumberISO(d, string, i) {
      var n = numberRe.exec(string.slice(i, i + 2));
      return n ? (d.V = +n[0], i + n[0].length) : -1;
    }

    function parseWeekNumberMonday(d, string, i) {
      var n = numberRe.exec(string.slice(i, i + 2));
      return n ? (d.W = +n[0], i + n[0].length) : -1;
    }

    function parseFullYear(d, string, i) {
      var n = numberRe.exec(string.slice(i, i + 4));
      return n ? (d.y = +n[0], i + n[0].length) : -1;
    }

    function parseYear(d, string, i) {
      var n = numberRe.exec(string.slice(i, i + 2));
      return n ? (d.y = +n[0] + (+n[0] > 68 ? 1900 : 2000), i + n[0].length) : -1;
    }

    function parseZone(d, string, i) {
      var n = /^(Z)|([+-]\d\d)(?::?(\d\d))?/.exec(string.slice(i, i + 6));
      return n ? (d.Z = n[1] ? 0 : -(n[2] + (n[3] || "00")), i + n[0].length) : -1;
    }

    function parseQuarter(d, string, i) {
      var n = numberRe.exec(string.slice(i, i + 1));
      return n ? (d.q = n[0] * 3 - 3, i + n[0].length) : -1;
    }

    function parseMonthNumber(d, string, i) {
      var n = numberRe.exec(string.slice(i, i + 2));
      return n ? (d.m = n[0] - 1, i + n[0].length) : -1;
    }

    function parseDayOfMonth(d, string, i) {
      var n = numberRe.exec(string.slice(i, i + 2));
      return n ? (d.d = +n[0], i + n[0].length) : -1;
    }

    function parseDayOfYear(d, string, i) {
      var n = numberRe.exec(string.slice(i, i + 3));
      return n ? (d.m = 0, d.d = +n[0], i + n[0].length) : -1;
    }

    function parseHour24(d, string, i) {
      var n = numberRe.exec(string.slice(i, i + 2));
      return n ? (d.H = +n[0], i + n[0].length) : -1;
    }

    function parseMinutes(d, string, i) {
      var n = numberRe.exec(string.slice(i, i + 2));
      return n ? (d.M = +n[0], i + n[0].length) : -1;
    }

    function parseSeconds(d, string, i) {
      var n = numberRe.exec(string.slice(i, i + 2));
      return n ? (d.S = +n[0], i + n[0].length) : -1;
    }

    function parseMilliseconds(d, string, i) {
      var n = numberRe.exec(string.slice(i, i + 3));
      return n ? (d.L = +n[0], i + n[0].length) : -1;
    }

    function parseMicroseconds(d, string, i) {
      var n = numberRe.exec(string.slice(i, i + 6));
      return n ? (d.L = Math.floor(n[0] / 1000), i + n[0].length) : -1;
    }

    function parseLiteralPercent(d, string, i) {
      var n = percentRe.exec(string.slice(i, i + 1));
      return n ? i + n[0].length : -1;
    }

    function parseUnixTimestamp(d, string, i) {
      var n = numberRe.exec(string.slice(i));
      return n ? (d.Q = +n[0], i + n[0].length) : -1;
    }

    function parseUnixTimestampSeconds(d, string, i) {
      var n = numberRe.exec(string.slice(i));
      return n ? (d.s = +n[0], i + n[0].length) : -1;
    }

    function formatDayOfMonth(d, p) {
      return pad(d.getDate(), p, 2);
    }

    function formatHour24(d, p) {
      return pad(d.getHours(), p, 2);
    }

    function formatHour12(d, p) {
      return pad(d.getHours() % 12 || 12, p, 2);
    }

    function formatDayOfYear(d, p) {
      return pad(1 + day.count(year(d), d), p, 3);
    }

    function formatMilliseconds(d, p) {
      return pad(d.getMilliseconds(), p, 3);
    }

    function formatMicroseconds(d, p) {
      return formatMilliseconds(d, p) + "000";
    }

    function formatMonthNumber(d, p) {
      return pad(d.getMonth() + 1, p, 2);
    }

    function formatMinutes(d, p) {
      return pad(d.getMinutes(), p, 2);
    }

    function formatSeconds(d, p) {
      return pad(d.getSeconds(), p, 2);
    }

    function formatWeekdayNumberMonday(d) {
      var day = d.getDay();
      return day === 0 ? 7 : day;
    }

    function formatWeekNumberSunday(d, p) {
      return pad(sunday.count(year(d) - 1, d), p, 2);
    }

    function formatWeekNumberISO(d, p) {
      var day = d.getDay();
      d = (day >= 4 || day === 0) ? thursday(d) : thursday.ceil(d);
      return pad(thursday.count(year(d), d) + (year(d).getDay() === 4), p, 2);
    }

    function formatWeekdayNumberSunday(d) {
      return d.getDay();
    }

    function formatWeekNumberMonday(d, p) {
      return pad(monday.count(year(d) - 1, d), p, 2);
    }

    function formatYear(d, p) {
      return pad(d.getFullYear() % 100, p, 2);
    }

    function formatFullYear(d, p) {
      return pad(d.getFullYear() % 10000, p, 4);
    }

    function formatZone(d) {
      var z = d.getTimezoneOffset();
      return (z > 0 ? "-" : (z *= -1, "+"))
          + pad(z / 60 | 0, "0", 2)
          + pad(z % 60, "0", 2);
    }

    function formatUTCDayOfMonth(d, p) {
      return pad(d.getUTCDate(), p, 2);
    }

    function formatUTCHour24(d, p) {
      return pad(d.getUTCHours(), p, 2);
    }

    function formatUTCHour12(d, p) {
      return pad(d.getUTCHours() % 12 || 12, p, 2);
    }

    function formatUTCDayOfYear(d, p) {
      return pad(1 + utcDay.count(utcYear(d), d), p, 3);
    }

    function formatUTCMilliseconds(d, p) {
      return pad(d.getUTCMilliseconds(), p, 3);
    }

    function formatUTCMicroseconds(d, p) {
      return formatUTCMilliseconds(d, p) + "000";
    }

    function formatUTCMonthNumber(d, p) {
      return pad(d.getUTCMonth() + 1, p, 2);
    }

    function formatUTCMinutes(d, p) {
      return pad(d.getUTCMinutes(), p, 2);
    }

    function formatUTCSeconds(d, p) {
      return pad(d.getUTCSeconds(), p, 2);
    }

    function formatUTCWeekdayNumberMonday(d) {
      var dow = d.getUTCDay();
      return dow === 0 ? 7 : dow;
    }

    function formatUTCWeekNumberSunday(d, p) {
      return pad(utcSunday.count(utcYear(d) - 1, d), p, 2);
    }

    function formatUTCWeekNumberISO(d, p) {
      var day = d.getUTCDay();
      d = (day >= 4 || day === 0) ? utcThursday(d) : utcThursday.ceil(d);
      return pad(utcThursday.count(utcYear(d), d) + (utcYear(d).getUTCDay() === 4), p, 2);
    }

    function formatUTCWeekdayNumberSunday(d) {
      return d.getUTCDay();
    }

    function formatUTCWeekNumberMonday(d, p) {
      return pad(utcMonday.count(utcYear(d) - 1, d), p, 2);
    }

    function formatUTCYear(d, p) {
      return pad(d.getUTCFullYear() % 100, p, 2);
    }

    function formatUTCFullYear(d, p) {
      return pad(d.getUTCFullYear() % 10000, p, 4);
    }

    function formatUTCZone() {
      return "+0000";
    }

    function formatLiteralPercent() {
      return "%";
    }

    function formatUnixTimestamp(d) {
      return +d;
    }

    function formatUnixTimestampSeconds(d) {
      return Math.floor(+d / 1000);
    }

    var locale$1;
    var timeFormat;
    var timeParse;
    var utcFormat;
    var utcParse;

    defaultLocale$1({
      dateTime: "%x, %X",
      date: "%-m/%-d/%Y",
      time: "%-I:%M:%S %p",
      periods: ["AM", "PM"],
      days: ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"],
      shortDays: ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"],
      months: ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"],
      shortMonths: ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"]
    });

    function defaultLocale$1(definition) {
      locale$1 = formatLocale$1(definition);
      timeFormat = locale$1.format;
      timeParse = locale$1.parse;
      utcFormat = locale$1.utcFormat;
      utcParse = locale$1.utcParse;
      return locale$1;
    }

    var isoSpecifier = "%Y-%m-%dT%H:%M:%S.%LZ";

    function formatIsoNative(date) {
      return date.toISOString();
    }

    var formatIso = Date.prototype.toISOString
        ? formatIsoNative
        : utcFormat(isoSpecifier);

    function parseIsoNative(string) {
      var date = new Date(string);
      return isNaN(date) ? null : date;
    }

    var parseIso = +new Date("2000-01-01T00:00:00.000Z")
        ? parseIsoNative
        : utcParse(isoSpecifier);

    const  points = [
    	{ x: 1979, y: 7.19 },
    	{ x: 1980, y: 7.83 },
    	{ x: 1981, y: 7.24 },
    	{ x: 1982, y: 7.44 },
    	{ x: 1983, y: 7.51 },
    	{ x: 1984, y: 7.10 },
    	{ x: 1985, y: 6.91 },
    	{ x: 1986, y: 7.53 },
    	{ x: 1987, y: 7.47 },
    	{ x: 1988, y: 7.48 },
    	{ x: 1989, y: 7.03 },
    	{ x: 1990, y: 6.23 },
    	{ x: 1991, y: 6.54 },
    	{ x: 1992, y: 7.54 },
    	{ x: 1993, y: 6.50 },
    	{ x: 1994, y: 7.18 },
    	{ x: 1995, y: 6.12 },
    	{ x: 1996, y: 7.87 },
    	{ x: 1997, y: 6.73 },
    	{ x: 1998, y: 6.55 },
    	{ x: 1999, y: 6.23 },
    	{ x: 2000, y: 6.31 },
    	{ x: 2001, y: 6.74 },
    	{ x: 2002, y: 5.95 },
    	{ x: 2003, y: 6.13 },
    	{ x: 2004, y: 6.04 },
    	{ x: 2005, y: 5.56 },
    	{ x: 2006, y: 5.91 },
    	{ x: 2007, y: 4.29 },
    	{ x: 2008, y: 4.72 },
    	{ x: 2009, y: 5.38 },
    	{ x: 2010, y: 4.92 },
    	{ x: 2011, y: 4.61 },
    	{ x: 2012, y: 3.62 },
    	{ x: 2013, y: 5.35 },
    	{ x: 2014, y: 5.28 },
    	{ x: 2015, y: 4.63 },
    	{ x: 2016, y: 4.72 }
    ];

    /* src/components/panels/PanelNetworkHashrate.svelte generated by Svelte v3.14.0 */
    const file$d = "src/components/panels/PanelNetworkHashrate.svelte";

    function get_each_context$3(ctx, list, i) {
    	const child_ctx = Object.create(ctx);
    	child_ctx.tick = list[i];
    	return child_ctx;
    }

    function get_each_context_1$2(ctx, list, i) {
    	const child_ctx = Object.create(ctx);
    	child_ctx.tick = list[i];
    	return child_ctx;
    }

    // (35:3) {#each yTicks as tick}
    function create_each_block_1$2(ctx) {
    	let g;
    	let line;
    	let text_1;
    	let t0_value = ctx.tick + "";
    	let t0;
    	let t1_value = (ctx.tick === 8 ? " million sq km" : "") + "";
    	let t1;
    	let g_class_value;
    	let g_transform_value;

    	const block = {
    		c: function create() {
    			g = svg_element("g");
    			line = svg_element("line");
    			text_1 = svg_element("text");
    			t0 = text(t0_value);
    			t1 = text(t1_value);
    			attr_dev(line, "x2", "100%");
    			attr_dev(line, "class", "svelte-lofqw8");
    			add_location(line, file$d, 36, 5, 1166);
    			attr_dev(text_1, "y", "-4");
    			attr_dev(text_1, "class", "svelte-lofqw8");
    			add_location(text_1, file$d, 37, 5, 1195);
    			attr_dev(g, "class", g_class_value = "tick tick-" + ctx.tick + " svelte-lofqw8");
    			attr_dev(g, "transform", g_transform_value = "translate(0, " + (ctx.yScale(ctx.tick) - ctx.padding.bottom) + ")");
    			add_location(g, file$d, 35, 4, 1074);
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, g, anchor);
    			append_dev(g, line);
    			append_dev(g, text_1);
    			append_dev(text_1, t0);
    			append_dev(text_1, t1);
    		},
    		p: function update(changed, ctx) {
    			if (changed.yScale && g_transform_value !== (g_transform_value = "translate(0, " + (ctx.yScale(ctx.tick) - ctx.padding.bottom) + ")")) {
    				attr_dev(g, "transform", g_transform_value);
    			}
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(g);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_each_block_1$2.name,
    		type: "each",
    		source: "(35:3) {#each yTicks as tick}",
    		ctx
    	});

    	return block;
    }

    // (45:3) {#each xTicks as tick}
    function create_each_block$3(ctx) {
    	let g;
    	let line;
    	let line_y__value;
    	let line_y__value_1;
    	let text_1;
    	let t_value = (ctx.width > 380 ? ctx.tick : formatMobile(ctx.tick)) + "";
    	let t;
    	let g_class_value;
    	let g_transform_value;

    	const block = {
    		c: function create() {
    			g = svg_element("g");
    			line = svg_element("line");
    			text_1 = svg_element("text");
    			t = text(t_value);
    			attr_dev(line, "y1", line_y__value = "-" + ctx.height);
    			attr_dev(line, "y2", line_y__value_1 = "-" + ctx.padding.bottom);
    			attr_dev(line, "x1", "0");
    			attr_dev(line, "x2", "0");
    			attr_dev(line, "class", "svelte-lofqw8");
    			add_location(line, file$d, 46, 5, 1444);
    			attr_dev(text_1, "y", "-2");
    			attr_dev(text_1, "class", "svelte-lofqw8");
    			add_location(text_1, file$d, 47, 5, 1515);
    			attr_dev(g, "class", g_class_value = "tick tick-" + ctx.tick + " svelte-lofqw8");
    			attr_dev(g, "transform", g_transform_value = "translate(" + ctx.xScale(ctx.tick) + "," + ctx.height + ")");
    			add_location(g, file$d, 45, 4, 1361);
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, g, anchor);
    			append_dev(g, line);
    			append_dev(g, text_1);
    			append_dev(text_1, t);
    		},
    		p: function update(changed, ctx) {
    			if (changed.height && line_y__value !== (line_y__value = "-" + ctx.height)) {
    				attr_dev(line, "y1", line_y__value);
    			}

    			if (changed.width && t_value !== (t_value = (ctx.width > 380 ? ctx.tick : formatMobile(ctx.tick)) + "")) set_data_dev(t, t_value);

    			if ((changed.xScale || changed.height) && g_transform_value !== (g_transform_value = "translate(" + ctx.xScale(ctx.tick) + "," + ctx.height + ")")) {
    				attr_dev(g, "transform", g_transform_value);
    			}
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(g);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_each_block$3.name,
    		type: "each",
    		source: "(45:3) {#each xTicks as tick}",
    		ctx
    	});

    	return block;
    }

    function create_fragment$d(ctx) {
    	let div1;
    	let div0;
    	let svg;
    	let g0;
    	let g0_transform_value;
    	let g1;
    	let path0;
    	let path1;
    	let div0_resize_listener;
    	let each_value_1 = ctx.yTicks;
    	let each_blocks_1 = [];

    	for (let i = 0; i < each_value_1.length; i += 1) {
    		each_blocks_1[i] = create_each_block_1$2(get_each_context_1$2(ctx, each_value_1, i));
    	}

    	let each_value = ctx.xTicks;
    	let each_blocks = [];

    	for (let i = 0; i < each_value.length; i += 1) {
    		each_blocks[i] = create_each_block$3(get_each_context$3(ctx, each_value, i));
    	}

    	const block = {
    		c: function create() {
    			div1 = element("div");
    			div0 = element("div");
    			svg = svg_element("svg");
    			g0 = svg_element("g");

    			for (let i = 0; i < each_blocks_1.length; i += 1) {
    				each_blocks_1[i].c();
    			}

    			g1 = svg_element("g");

    			for (let i = 0; i < each_blocks.length; i += 1) {
    				each_blocks[i].c();
    			}

    			path0 = svg_element("path");
    			path1 = svg_element("path");
    			attr_dev(g0, "class", "axis y-axis");
    			attr_dev(g0, "transform", g0_transform_value = "translate(0, " + ctx.padding.top + ")");
    			add_location(g0, file$d, 33, 2, 980);
    			attr_dev(g1, "class", "axis x-axis svelte-lofqw8");
    			add_location(g1, file$d, 43, 2, 1307);
    			attr_dev(path0, "class", "path-area svelte-lofqw8");
    			attr_dev(path0, "d", ctx.area);
    			add_location(path0, file$d, 53, 2, 1623);
    			attr_dev(path1, "class", "path-line svelte-lofqw8");
    			attr_dev(path1, "d", ctx.path);
    			add_location(path1, file$d, 54, 2, 1666);
    			attr_dev(svg, "class", "rwrap svelte-lofqw8");
    			add_location(svg, file$d, 31, 1, 940);
    			attr_dev(div0, "class", "chart svelte-lofqw8");
    			add_render_callback(() => ctx.div0_resize_handler.call(div0));
    			add_location(div0, file$d, 30, 0, 867);
    			attr_dev(div1, "class", "rwrap flx");
    			add_location(div1, file$d, 29, 0, 843);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div1, anchor);
    			append_dev(div1, div0);
    			append_dev(div0, svg);
    			append_dev(svg, g0);

    			for (let i = 0; i < each_blocks_1.length; i += 1) {
    				each_blocks_1[i].m(g0, null);
    			}

    			append_dev(svg, g1);

    			for (let i = 0; i < each_blocks.length; i += 1) {
    				each_blocks[i].m(g1, null);
    			}

    			append_dev(svg, path0);
    			append_dev(svg, path1);
    			div0_resize_listener = add_resize_listener(div0, ctx.div0_resize_handler.bind(div0));
    		},
    		p: function update(changed, ctx) {
    			if (changed.yTicks || changed.yScale || changed.padding) {
    				each_value_1 = ctx.yTicks;
    				let i;

    				for (i = 0; i < each_value_1.length; i += 1) {
    					const child_ctx = get_each_context_1$2(ctx, each_value_1, i);

    					if (each_blocks_1[i]) {
    						each_blocks_1[i].p(changed, child_ctx);
    					} else {
    						each_blocks_1[i] = create_each_block_1$2(child_ctx);
    						each_blocks_1[i].c();
    						each_blocks_1[i].m(g0, null);
    					}
    				}

    				for (; i < each_blocks_1.length; i += 1) {
    					each_blocks_1[i].d(1);
    				}

    				each_blocks_1.length = each_value_1.length;
    			}

    			if (changed.xTicks || changed.xScale || changed.height || changed.width || changed.formatMobile || changed.padding) {
    				each_value = ctx.xTicks;
    				let i;

    				for (i = 0; i < each_value.length; i += 1) {
    					const child_ctx = get_each_context$3(ctx, each_value, i);

    					if (each_blocks[i]) {
    						each_blocks[i].p(changed, child_ctx);
    					} else {
    						each_blocks[i] = create_each_block$3(child_ctx);
    						each_blocks[i].c();
    						each_blocks[i].m(g1, null);
    					}
    				}

    				for (; i < each_blocks.length; i += 1) {
    					each_blocks[i].d(1);
    				}

    				each_blocks.length = each_value.length;
    			}

    			if (changed.area) {
    				attr_dev(path0, "d", ctx.area);
    			}

    			if (changed.path) {
    				attr_dev(path1, "d", ctx.path);
    			}
    		},
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div1);
    			destroy_each(each_blocks_1, detaching);
    			destroy_each(each_blocks, detaching);
    			div0_resize_listener.cancel();
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$d.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    function formatMobile(tick) {
    	return "'" + tick % 100;
    }

    function instance$9($$self, $$props, $$invalidate) {
    	const yTicks = [0, 2, 4, 6, 8];
    	const xTicks = [1980, 1990, 2000, 2010];
    	const padding = { top: 20, right: 15, bottom: 20, left: 25 };
    	let width = 500;
    	let height = 200;

    	function div0_resize_handler() {
    		width = this.clientWidth;
    		height = this.clientHeight;
    		$$invalidate("width", width);
    		$$invalidate("height", height);
    	}

    	$$self.$capture_state = () => {
    		return {};
    	};

    	$$self.$inject_state = $$props => {
    		if ("width" in $$props) $$invalidate("width", width = $$props.width);
    		if ("height" in $$props) $$invalidate("height", height = $$props.height);
    		if ("xScale" in $$props) $$invalidate("xScale", xScale = $$props.xScale);
    		if ("minX" in $$props) $$invalidate("minX", minX = $$props.minX);
    		if ("maxX" in $$props) $$invalidate("maxX", maxX = $$props.maxX);
    		if ("yScale" in $$props) $$invalidate("yScale", yScale = $$props.yScale);
    		if ("path" in $$props) $$invalidate("path", path = $$props.path);
    		if ("area" in $$props) $$invalidate("area", area = $$props.area);
    	};

    	let xScale;
    	let yScale;
    	let minX;
    	let maxX;
    	let path;
    	let area;

    	$$self.$$.update = (changed = { minX: 1, maxX: 1, width: 1, height: 1, xScale: 1, yScale: 1, path: 1 }) => {
    		if (changed.minX || changed.maxX || changed.width) {
    			 $$invalidate("xScale", xScale = linear$1().domain([minX, maxX]).range([padding.left, width - padding.right]));
    		}

    		if (changed.height) {
    			 $$invalidate("yScale", yScale = linear$1().domain([Math.min.apply(null, yTicks), Math.max.apply(null, yTicks)]).range([height - padding.bottom, padding.top]));
    		}

    		if (changed.xScale || changed.yScale) {
    			 $$invalidate("path", path = `M${points.map(p => `${xScale(p.x)},${yScale(p.y)}`).join("L")}`);
    		}

    		if (changed.path || changed.xScale || changed.maxX || changed.yScale || changed.minX) {
    			 $$invalidate("area", area = `${path}L${xScale(maxX)},${yScale(0)}L${xScale(minX)},${yScale(0)}Z`);
    		}
    	};

    	 $$invalidate("minX", minX = points[0].x);
    	 $$invalidate("maxX", maxX = points[points.length - 1].x);

    	return {
    		yTicks,
    		xTicks,
    		padding,
    		width,
    		height,
    		xScale,
    		yScale,
    		path,
    		area,
    		div0_resize_handler
    	};
    }

    class PanelNetworkHashrate extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, instance$9, create_fragment$d, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PanelNetworkHashrate",
    			options,
    			id: create_fragment$d.name
    		});
    	}
    }

    /* src/components/panels/PanelLocalHashrate.svelte generated by Svelte v3.14.0 */
    const file$e = "src/components/panels/PanelLocalHashrate.svelte";

    function get_each_context$4(ctx, list, i) {
    	const child_ctx = Object.create(ctx);
    	child_ctx.tick = list[i];
    	return child_ctx;
    }

    function get_each_context_1$3(ctx, list, i) {
    	const child_ctx = Object.create(ctx);
    	child_ctx.tick = list[i];
    	return child_ctx;
    }

    // (35:3) {#each yTicks as tick}
    function create_each_block_1$3(ctx) {
    	let g;
    	let line;
    	let text_1;
    	let t0_value = ctx.tick + "";
    	let t0;
    	let t1_value = (ctx.tick === 8 ? " million sq km" : "") + "";
    	let t1;
    	let g_class_value;
    	let g_transform_value;

    	const block = {
    		c: function create() {
    			g = svg_element("g");
    			line = svg_element("line");
    			text_1 = svg_element("text");
    			t0 = text(t0_value);
    			t1 = text(t1_value);
    			attr_dev(line, "x2", "100%");
    			attr_dev(line, "class", "svelte-1vgh4lf");
    			add_location(line, file$e, 36, 5, 1166);
    			attr_dev(text_1, "y", "-4");
    			attr_dev(text_1, "class", "svelte-1vgh4lf");
    			add_location(text_1, file$e, 37, 5, 1195);
    			attr_dev(g, "class", g_class_value = "tick tick-" + ctx.tick + " svelte-1vgh4lf");
    			attr_dev(g, "transform", g_transform_value = "translate(0, " + (ctx.yScale(ctx.tick) - ctx.padding.bottom) + ")");
    			add_location(g, file$e, 35, 4, 1074);
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, g, anchor);
    			append_dev(g, line);
    			append_dev(g, text_1);
    			append_dev(text_1, t0);
    			append_dev(text_1, t1);
    		},
    		p: function update(changed, ctx) {
    			if (changed.yScale && g_transform_value !== (g_transform_value = "translate(0, " + (ctx.yScale(ctx.tick) - ctx.padding.bottom) + ")")) {
    				attr_dev(g, "transform", g_transform_value);
    			}
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(g);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_each_block_1$3.name,
    		type: "each",
    		source: "(35:3) {#each yTicks as tick}",
    		ctx
    	});

    	return block;
    }

    // (45:3) {#each xTicks as tick}
    function create_each_block$4(ctx) {
    	let g;
    	let line;
    	let line_y__value;
    	let line_y__value_1;
    	let text_1;
    	let t_value = (ctx.width > 380 ? ctx.tick : formatMobile$1(ctx.tick)) + "";
    	let t;
    	let g_class_value;
    	let g_transform_value;

    	const block = {
    		c: function create() {
    			g = svg_element("g");
    			line = svg_element("line");
    			text_1 = svg_element("text");
    			t = text(t_value);
    			attr_dev(line, "y1", line_y__value = "-" + ctx.height);
    			attr_dev(line, "y2", line_y__value_1 = "-" + ctx.padding.bottom);
    			attr_dev(line, "x1", "0");
    			attr_dev(line, "x2", "0");
    			attr_dev(line, "class", "svelte-1vgh4lf");
    			add_location(line, file$e, 46, 5, 1444);
    			attr_dev(text_1, "y", "-2");
    			attr_dev(text_1, "class", "svelte-1vgh4lf");
    			add_location(text_1, file$e, 47, 5, 1515);
    			attr_dev(g, "class", g_class_value = "tick tick-" + ctx.tick + " svelte-1vgh4lf");
    			attr_dev(g, "transform", g_transform_value = "translate(" + ctx.xScale(ctx.tick) + "," + ctx.height + ")");
    			add_location(g, file$e, 45, 4, 1361);
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, g, anchor);
    			append_dev(g, line);
    			append_dev(g, text_1);
    			append_dev(text_1, t);
    		},
    		p: function update(changed, ctx) {
    			if (changed.height && line_y__value !== (line_y__value = "-" + ctx.height)) {
    				attr_dev(line, "y1", line_y__value);
    			}

    			if (changed.width && t_value !== (t_value = (ctx.width > 380 ? ctx.tick : formatMobile$1(ctx.tick)) + "")) set_data_dev(t, t_value);

    			if ((changed.xScale || changed.height) && g_transform_value !== (g_transform_value = "translate(" + ctx.xScale(ctx.tick) + "," + ctx.height + ")")) {
    				attr_dev(g, "transform", g_transform_value);
    			}
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(g);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_each_block$4.name,
    		type: "each",
    		source: "(45:3) {#each xTicks as tick}",
    		ctx
    	});

    	return block;
    }

    function create_fragment$e(ctx) {
    	let div1;
    	let div0;
    	let svg;
    	let g0;
    	let g0_transform_value;
    	let g1;
    	let path0;
    	let path1;
    	let div0_resize_listener;
    	let each_value_1 = ctx.yTicks;
    	let each_blocks_1 = [];

    	for (let i = 0; i < each_value_1.length; i += 1) {
    		each_blocks_1[i] = create_each_block_1$3(get_each_context_1$3(ctx, each_value_1, i));
    	}

    	let each_value = ctx.xTicks;
    	let each_blocks = [];

    	for (let i = 0; i < each_value.length; i += 1) {
    		each_blocks[i] = create_each_block$4(get_each_context$4(ctx, each_value, i));
    	}

    	const block = {
    		c: function create() {
    			div1 = element("div");
    			div0 = element("div");
    			svg = svg_element("svg");
    			g0 = svg_element("g");

    			for (let i = 0; i < each_blocks_1.length; i += 1) {
    				each_blocks_1[i].c();
    			}

    			g1 = svg_element("g");

    			for (let i = 0; i < each_blocks.length; i += 1) {
    				each_blocks[i].c();
    			}

    			path0 = svg_element("path");
    			path1 = svg_element("path");
    			attr_dev(g0, "class", "axis y-axis");
    			attr_dev(g0, "transform", g0_transform_value = "translate(0, " + ctx.padding.top + ")");
    			add_location(g0, file$e, 33, 2, 980);
    			attr_dev(g1, "class", "axis x-axis svelte-1vgh4lf");
    			add_location(g1, file$e, 43, 2, 1307);
    			attr_dev(path0, "class", "path-area svelte-1vgh4lf");
    			attr_dev(path0, "d", ctx.area);
    			add_location(path0, file$e, 53, 2, 1623);
    			attr_dev(path1, "class", "path-line svelte-1vgh4lf");
    			attr_dev(path1, "d", ctx.path);
    			add_location(path1, file$e, 54, 2, 1666);
    			attr_dev(svg, "class", "rwrap svelte-1vgh4lf");
    			add_location(svg, file$e, 31, 1, 940);
    			attr_dev(div0, "class", "chart svelte-1vgh4lf");
    			add_render_callback(() => ctx.div0_resize_handler.call(div0));
    			add_location(div0, file$e, 30, 0, 867);
    			attr_dev(div1, "class", "rwrap flx");
    			add_location(div1, file$e, 29, 0, 843);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div1, anchor);
    			append_dev(div1, div0);
    			append_dev(div0, svg);
    			append_dev(svg, g0);

    			for (let i = 0; i < each_blocks_1.length; i += 1) {
    				each_blocks_1[i].m(g0, null);
    			}

    			append_dev(svg, g1);

    			for (let i = 0; i < each_blocks.length; i += 1) {
    				each_blocks[i].m(g1, null);
    			}

    			append_dev(svg, path0);
    			append_dev(svg, path1);
    			div0_resize_listener = add_resize_listener(div0, ctx.div0_resize_handler.bind(div0));
    		},
    		p: function update(changed, ctx) {
    			if (changed.yTicks || changed.yScale || changed.padding) {
    				each_value_1 = ctx.yTicks;
    				let i;

    				for (i = 0; i < each_value_1.length; i += 1) {
    					const child_ctx = get_each_context_1$3(ctx, each_value_1, i);

    					if (each_blocks_1[i]) {
    						each_blocks_1[i].p(changed, child_ctx);
    					} else {
    						each_blocks_1[i] = create_each_block_1$3(child_ctx);
    						each_blocks_1[i].c();
    						each_blocks_1[i].m(g0, null);
    					}
    				}

    				for (; i < each_blocks_1.length; i += 1) {
    					each_blocks_1[i].d(1);
    				}

    				each_blocks_1.length = each_value_1.length;
    			}

    			if (changed.xTicks || changed.xScale || changed.height || changed.width || changed.formatMobile || changed.padding) {
    				each_value = ctx.xTicks;
    				let i;

    				for (i = 0; i < each_value.length; i += 1) {
    					const child_ctx = get_each_context$4(ctx, each_value, i);

    					if (each_blocks[i]) {
    						each_blocks[i].p(changed, child_ctx);
    					} else {
    						each_blocks[i] = create_each_block$4(child_ctx);
    						each_blocks[i].c();
    						each_blocks[i].m(g1, null);
    					}
    				}

    				for (; i < each_blocks.length; i += 1) {
    					each_blocks[i].d(1);
    				}

    				each_blocks.length = each_value.length;
    			}

    			if (changed.area) {
    				attr_dev(path0, "d", ctx.area);
    			}

    			if (changed.path) {
    				attr_dev(path1, "d", ctx.path);
    			}
    		},
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div1);
    			destroy_each(each_blocks_1, detaching);
    			destroy_each(each_blocks, detaching);
    			div0_resize_listener.cancel();
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$e.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    function formatMobile$1(tick) {
    	return "'" + tick % 100;
    }

    function instance$a($$self, $$props, $$invalidate) {
    	const yTicks = [0, 2, 4, 6, 8];
    	const xTicks = [1980, 1990, 2000, 2010];
    	const padding = { top: 20, right: 15, bottom: 20, left: 25 };
    	let width = 500;
    	let height = 200;

    	function div0_resize_handler() {
    		width = this.clientWidth;
    		height = this.clientHeight;
    		$$invalidate("width", width);
    		$$invalidate("height", height);
    	}

    	$$self.$capture_state = () => {
    		return {};
    	};

    	$$self.$inject_state = $$props => {
    		if ("width" in $$props) $$invalidate("width", width = $$props.width);
    		if ("height" in $$props) $$invalidate("height", height = $$props.height);
    		if ("xScale" in $$props) $$invalidate("xScale", xScale = $$props.xScale);
    		if ("minX" in $$props) $$invalidate("minX", minX = $$props.minX);
    		if ("maxX" in $$props) $$invalidate("maxX", maxX = $$props.maxX);
    		if ("yScale" in $$props) $$invalidate("yScale", yScale = $$props.yScale);
    		if ("path" in $$props) $$invalidate("path", path = $$props.path);
    		if ("area" in $$props) $$invalidate("area", area = $$props.area);
    	};

    	let xScale;
    	let yScale;
    	let minX;
    	let maxX;
    	let path;
    	let area;

    	$$self.$$.update = (changed = { minX: 1, maxX: 1, width: 1, height: 1, xScale: 1, yScale: 1, path: 1 }) => {
    		if (changed.minX || changed.maxX || changed.width) {
    			 $$invalidate("xScale", xScale = linear$1().domain([minX, maxX]).range([padding.left, width - padding.right]));
    		}

    		if (changed.height) {
    			 $$invalidate("yScale", yScale = linear$1().domain([Math.min.apply(null, yTicks), Math.max.apply(null, yTicks)]).range([height - padding.bottom, padding.top]));
    		}

    		if (changed.xScale || changed.yScale) {
    			 $$invalidate("path", path = `M${points.map(p => `${xScale(p.x)},${yScale(p.y)}`).join("L")}`);
    		}

    		if (changed.path || changed.xScale || changed.maxX || changed.yScale || changed.minX) {
    			 $$invalidate("area", area = `${path}L${xScale(maxX)},${yScale(0)}L${xScale(minX)},${yScale(0)}Z`);
    		}
    	};

    	 $$invalidate("minX", minX = points[0].x);
    	 $$invalidate("maxX", maxX = points[points.length - 1].x);

    	return {
    		yTicks,
    		xTicks,
    		padding,
    		width,
    		height,
    		xScale,
    		yScale,
    		path,
    		area,
    		div0_resize_handler
    	};
    }

    class PanelLocalHashrate extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, instance$a, create_fragment$e, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PanelLocalHashrate",
    			options,
    			id: create_fragment$e.name
    		});
    	}
    }

    /* src/components/panels/PanelStatus.svelte generated by Svelte v3.14.0 */

    const file$f = "src/components/panels/PanelStatus.svelte";

    function create_fragment$f(ctx) {
    	let div;
    	let ul;
    	let li0;
    	let span0;
    	let t0;
    	let span1;
    	let t2;
    	let strong0;
    	let span2;
    	let t3;
    	let li1;
    	let span3;
    	let t4;
    	let span4;
    	let t6;
    	let strong1;
    	let span5;
    	let t7;
    	let li2;
    	let span6;
    	let t8;
    	let span7;
    	let t10;
    	let strong2;
    	let span8;
    	let t11;
    	let li3;
    	let span9;
    	let t12;
    	let span10;
    	let t14;
    	let strong3;
    	let span11;
    	let t15;
    	let li4;
    	let span12;
    	let t16;
    	let span13;
    	let t18;
    	let strong4;
    	let span14;
    	let t19;
    	let li5;
    	let span15;
    	let t20;
    	let span16;
    	let t22;
    	let strong5;
    	let span17;
    	let t23;
    	let li6;
    	let span18;
    	let t24;
    	let span19;
    	let t26;
    	let strong6;
    	let span20;
    	let t27;
    	let li7;
    	let span21;
    	let t28;
    	let span22;
    	let t30;
    	let strong7;
    	let span23;

    	const block = {
    		c: function create() {
    			div = element("div");
    			ul = element("ul");
    			li0 = element("li");
    			span0 = element("span");
    			t0 = space();
    			span1 = element("span");
    			span1.textContent = "Version:";
    			t2 = space();
    			strong0 = element("strong");
    			span2 = element("span");
    			t3 = space();
    			li1 = element("li");
    			span3 = element("span");
    			t4 = space();
    			span4 = element("span");
    			span4.textContent = "Wallet version:";
    			t6 = space();
    			strong1 = element("strong");
    			span5 = element("span");
    			t7 = space();
    			li2 = element("li");
    			span6 = element("span");
    			t8 = space();
    			span7 = element("span");
    			span7.textContent = "Uptime:";
    			t10 = space();
    			strong2 = element("strong");
    			span8 = element("span");
    			t11 = space();
    			li3 = element("li");
    			span9 = element("span");
    			t12 = space();
    			span10 = element("span");
    			span10.textContent = "Memory:";
    			t14 = space();
    			strong3 = element("strong");
    			span11 = element("span");
    			t15 = space();
    			li4 = element("li");
    			span12 = element("span");
    			t16 = space();
    			span13 = element("span");
    			span13.textContent = "Disk:";
    			t18 = space();
    			strong4 = element("strong");
    			span14 = element("span");
    			t19 = space();
    			li5 = element("li");
    			span15 = element("span");
    			t20 = space();
    			span16 = element("span");
    			span16.textContent = "Chain:";
    			t22 = space();
    			strong5 = element("strong");
    			span17 = element("span");
    			t23 = space();
    			li6 = element("li");
    			span18 = element("span");
    			t24 = space();
    			span19 = element("span");
    			span19.textContent = "Blocks:";
    			t26 = space();
    			strong6 = element("strong");
    			span20 = element("span");
    			t27 = space();
    			li7 = element("li");
    			span21 = element("span");
    			t28 = space();
    			span22 = element("span");
    			span22.textContent = "Connections:";
    			t30 = space();
    			strong7 = element("strong");
    			span23 = element("span");
    			attr_dev(span0, "class", "rcx2");
    			add_location(span0, file$f, 3, 11, 140);
    			attr_dev(span1, "class", "rcx4");
    			add_location(span1, file$f, 4, 11, 178);
    			attr_dev(span2, "v-html", "this.duOSys.status.ver");
    			add_location(span2, file$f, 5, 32, 246);
    			attr_dev(strong0, "class", "rcx6");
    			add_location(strong0, file$f, 5, 11, 225);
    			attr_dev(li0, "class", "flx fwd spb htg rr");
    			add_location(li0, file$f, 2, 2, 97);
    			attr_dev(span3, "class", "rcx2");
    			add_location(span3, file$f, 8, 11, 370);
    			attr_dev(span4, "class", "rcx4");
    			add_location(span4, file$f, 9, 11, 408);
    			attr_dev(span5, "v-html", "this.duOSys.status.walletver.podjsonrpcapi.versionstring");
    			add_location(span5, file$f, 10, 32, 483);
    			attr_dev(strong1, "class", "rcx6");
    			add_location(strong1, file$f, 10, 11, 462);
    			attr_dev(li1, "class", "flx fwd spb htg rr");
    			add_location(li1, file$f, 7, 10, 327);
    			attr_dev(span6, "class", "rcx2");
    			add_location(span6, file$f, 13, 11, 641);
    			attr_dev(span7, "class", "rcx4");
    			add_location(span7, file$f, 14, 11, 679);
    			attr_dev(span8, "v-html", "this.duOSys.status.uptime");
    			add_location(span8, file$f, 15, 32, 746);
    			attr_dev(strong2, "class", "rcx6");
    			add_location(strong2, file$f, 15, 11, 725);
    			attr_dev(li2, "class", "flx fwd spb htg rr");
    			add_location(li2, file$f, 12, 10, 598);
    			attr_dev(span9, "class", "rcx2");
    			add_location(span9, file$f, 18, 11, 873);
    			attr_dev(span10, "class", "rcx4");
    			add_location(span10, file$f, 19, 11, 911);
    			attr_dev(span11, "v-html", "this.duOSys.status.net");
    			add_location(span11, file$f, 20, 32, 978);
    			attr_dev(strong3, "class", "rcx6");
    			add_location(strong3, file$f, 20, 11, 957);
    			attr_dev(li3, "class", "flx fwd spb htg rr");
    			add_location(li3, file$f, 17, 10, 830);
    			attr_dev(span12, "class", "rcx2");
    			add_location(span12, file$f, 24, 19, 1111);
    			attr_dev(span13, "class", "rcx4");
    			add_location(span13, file$f, 25, 19, 1157);
    			attr_dev(span14, "v-html", "this.duOSys.status.ver");
    			add_location(span14, file$f, 26, 40, 1230);
    			attr_dev(strong4, "class", "rcx6");
    			add_location(strong4, file$f, 26, 19, 1209);
    			attr_dev(li4, "class", "flx fwd spb htg rr");
    			add_location(li4, file$f, 23, 10, 1060);
    			attr_dev(span15, "class", "rcx2");
    			add_location(span15, file$f, 29, 19, 1378);
    			attr_dev(span16, "class", "rcx4");
    			add_location(span16, file$f, 30, 19, 1424);
    			attr_dev(span17, "v-html", "this.duOSys.status.walletver.podjsonrpcapi.versionstring");
    			add_location(span17, file$f, 31, 40, 1498);
    			attr_dev(strong5, "class", "rcx6");
    			add_location(strong5, file$f, 31, 19, 1477);
    			attr_dev(li5, "class", "flx fwd spb htg rr");
    			add_location(li5, file$f, 28, 18, 1327);
    			attr_dev(span18, "class", "rcx2");
    			add_location(span18, file$f, 34, 19, 1680);
    			attr_dev(span19, "class", "rcx4");
    			add_location(span19, file$f, 35, 19, 1726);
    			attr_dev(span20, "v-html", "this.duOSys.status.uptime");
    			add_location(span20, file$f, 36, 40, 1801);
    			attr_dev(strong6, "class", "rcx6");
    			add_location(strong6, file$f, 36, 19, 1780);
    			attr_dev(li6, "class", "flx fwd spb htg rr");
    			add_location(li6, file$f, 33, 18, 1629);
    			attr_dev(span21, "class", "rcx2");
    			add_location(span21, file$f, 39, 19, 1952);
    			attr_dev(span22, "class", "rcx4");
    			add_location(span22, file$f, 40, 19, 1998);
    			attr_dev(span23, "v-html", "this.duOSys.status.net");
    			add_location(span23, file$f, 41, 40, 2078);
    			attr_dev(strong7, "class", "rcx6");
    			add_location(strong7, file$f, 41, 19, 2057);
    			attr_dev(li7, "class", "flx fwd spb htg rr");
    			add_location(li7, file$f, 38, 18, 1901);
    			attr_dev(ul, "class", "rf flx flc noMargin noPadding justifyEvenly");
    			add_location(ul, file$f, 1, 1, 38);
    			attr_dev(div, "id", "panelstatus");
    			attr_dev(div, "class", "Info");
    			add_location(div, file$f, 0, 0, 0);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div, anchor);
    			append_dev(div, ul);
    			append_dev(ul, li0);
    			append_dev(li0, span0);
    			append_dev(li0, t0);
    			append_dev(li0, span1);
    			append_dev(li0, t2);
    			append_dev(li0, strong0);
    			append_dev(strong0, span2);
    			append_dev(ul, t3);
    			append_dev(ul, li1);
    			append_dev(li1, span3);
    			append_dev(li1, t4);
    			append_dev(li1, span4);
    			append_dev(li1, t6);
    			append_dev(li1, strong1);
    			append_dev(strong1, span5);
    			append_dev(ul, t7);
    			append_dev(ul, li2);
    			append_dev(li2, span6);
    			append_dev(li2, t8);
    			append_dev(li2, span7);
    			append_dev(li2, t10);
    			append_dev(li2, strong2);
    			append_dev(strong2, span8);
    			append_dev(ul, t11);
    			append_dev(ul, li3);
    			append_dev(li3, span9);
    			append_dev(li3, t12);
    			append_dev(li3, span10);
    			append_dev(li3, t14);
    			append_dev(li3, strong3);
    			append_dev(strong3, span11);
    			append_dev(ul, t15);
    			append_dev(ul, li4);
    			append_dev(li4, span12);
    			append_dev(li4, t16);
    			append_dev(li4, span13);
    			append_dev(li4, t18);
    			append_dev(li4, strong4);
    			append_dev(strong4, span14);
    			append_dev(ul, t19);
    			append_dev(ul, li5);
    			append_dev(li5, span15);
    			append_dev(li5, t20);
    			append_dev(li5, span16);
    			append_dev(li5, t22);
    			append_dev(li5, strong5);
    			append_dev(strong5, span17);
    			append_dev(ul, t23);
    			append_dev(ul, li6);
    			append_dev(li6, span18);
    			append_dev(li6, t24);
    			append_dev(li6, span19);
    			append_dev(li6, t26);
    			append_dev(li6, strong6);
    			append_dev(strong6, span20);
    			append_dev(ul, t27);
    			append_dev(ul, li7);
    			append_dev(li7, span21);
    			append_dev(li7, t28);
    			append_dev(li7, span22);
    			append_dev(li7, t30);
    			append_dev(li7, strong7);
    			append_dev(strong7, span23);
    		},
    		p: noop,
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$f.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class PanelStatus extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$f, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PanelStatus",
    			options,
    			id: create_fragment$f.name
    		});
    	}
    }

    const loading = writable(false);

    const lasTxs = readable([], function start(set) {
        const interval = setInterval(() => {
            fetch(`http://127.0.0.1:3999/lastxs`)
                .then(resp => resp.json())
                .then(data => {
                    set(data);
                });
        }, 100);
        return function stop() {
            clearInterval(interval);
        };
    });

    /* src/components/panels/PanelLatestTx.svelte generated by Svelte v3.14.0 */
    const file$g = "src/components/panels/PanelLatestTx.svelte";

    function create_fragment$g(ctx) {
    	let div;
    	let current;

    	const datatable = new DataTable({
    			props: {
    				lasTxs,
    				loading,
    				columns: [
    					{
    						label: "ID",
    						field: "id",
    						class: "md:w-10"
    					},
    					{
    						label: "Ep.",
    						value: func$1,
    						class: "md:w-10",
    						editable: false
    					},
    					{ field: "name", class: "md:w-10" },
    					{
    						field: "summary",
    						textarea: true,
    						value: func_1$1,
    						class: "text-sm text-gray-700 caption md:w-full sm:w-64"
    					},
    					{
    						field: "thumbnail",
    						value: func_2,
    						class: "w-48",
    						sortable: false,
    						editable: false
    					}
    				]
    			},
    			$$inline: true
    		});

    	datatable.$on("update", ctx.update_handler);

    	const block = {
    		c: function create() {
    			div = element("div");
    			create_component(datatable.$$.fragment);
    			attr_dev(div, "class", "rwrap flx");
    			add_location(div, file$g, 6, 0, 157);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div, anchor);
    			mount_component(datatable, div, null);
    			current = true;
    		},
    		p: function update(changed, ctx) {
    			const datatable_changes = {};
    			if (changed.lasTxs) datatable_changes.lasTxs = lasTxs;
    			datatable.$set(datatable_changes);
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(datatable.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(datatable.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div);
    			destroy_component(datatable);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$g.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    const func$1 = v => `S${v.season}E${v.number}`;
    const func_1$1 = v => v && v.summary ? v.summary : "";

    const func_2 = v => v && v.image
    ? `<img src="${v.image.medium.replace("http", "https")}" height="70" alt="${v.name}">`
    : "";

    function instance$b($$self) {
    	const update_handler = ({ detail }) => {
    		const { column, item, value } = detail;
    		const index = lasTxs.findIndex(i => i.id === item.id);
    		lasTxs[index][column.field] = value;
    	};

    	$$self.$capture_state = () => {
    		return {};
    	};

    	$$self.$inject_state = $$props => {
    		
    	};

    	return { update_handler };
    }

    class PanelLatestTx extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, instance$b, create_fragment$g, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PanelLatestTx",
    			options,
    			id: create_fragment$g.name
    		});
    	}
    }

    const status = readable([], function start(set) {
        const interval = setInterval(() => {
            fetch(`http://127.0.0.1:3999/status`)
                .then(resp => resp.json())
                .then(data => {
                    set(data);
                });
        }, 1000);
        return function stop() {
            clearInterval(interval);
        };
    });

    /* src/components/pages/PageOverview.svelte generated by Svelte v3.14.0 */
    const file$h = "src/components/pages/PageOverview.svelte";

    function create_fragment$h(ctx) {
    	let t0;
    	let main;
    	let div0;
    	let t1;
    	let div1;
    	let t2;
    	let div2;
    	let t3;
    	let div3;
    	let t4;
    	let div4;
    	let t5;
    	let div5;
    	let t6;
    	let div6;
    	let t7;
    	let div7;
    	let current;

    	const panelbalance = new PanelBalance({
    			props: { balance: status.balance },
    			$$inline: true
    		});

    	const panelsend = new PanelSend({ $$inline: true });

    	const panelnetworkhashrate = new PanelNetworkHashrate({
    			props: { hashrate: status.networkhashrate },
    			$$inline: true
    		});

    	const panellocalhashrate = new PanelLocalHashrate({
    			props: { hashrate: status.networkhashrate },
    			$$inline: true
    		});

    	const panelstatus = new PanelStatus({ props: { status }, $$inline: true });
    	const panellatesttx = new PanelLatestTx({ $$inline: true });

    	const block = {
    		c: function create() {
    			t0 = space();
    			main = element("main");
    			div0 = element("div");
    			create_component(panelbalance.$$.fragment);
    			t1 = space();
    			div1 = element("div");
    			create_component(panelsend.$$.fragment);
    			t2 = space();
    			div2 = element("div");
    			create_component(panelnetworkhashrate.$$.fragment);
    			t3 = space();
    			div3 = element("div");
    			create_component(panellocalhashrate.$$.fragment);
    			t4 = space();
    			div4 = element("div");
    			create_component(panelstatus.$$.fragment);
    			t5 = space();
    			div5 = element("div");
    			create_component(panellatesttx.$$.fragment);
    			t6 = space();
    			div6 = element("div");
    			t7 = space();
    			div7 = element("div");
    			document.title = "Overview";
    			attr_dev(div0, "id", "panelwalletstatus");
    			attr_dev(div0, "class", "Balance");
    			add_location(div0, file$h, 26, 8, 781);
    			attr_dev(div1, "id", "panelsend");
    			attr_dev(div1, "class", "Send");
    			add_location(div1, file$h, 29, 8, 902);
    			attr_dev(div2, "id", "panelnetworkhashrate");
    			attr_dev(div2, "class", "NetHash");
    			add_location(div2, file$h, 32, 8, 986);
    			attr_dev(div3, "id", "panellocalhashrate");
    			attr_dev(div3, "class", "LocalHash");
    			add_location(div3, file$h, 35, 8, 1128);
    			attr_dev(div4, "id", "panelstatus");
    			attr_dev(div4, "class", "Status");
    			add_location(div4, file$h, 38, 8, 1268);
    			attr_dev(div5, "id", "paneltxsex");
    			attr_dev(div5, "class", "Txs");
    			add_location(div5, file$h, 41, 8, 1373);
    			attr_dev(div6, "class", "Info");
    			add_location(div6, file$h, 44, 8, 1460);
    			attr_dev(div7, "class", "Time");
    			add_location(div7, file$h, 46, 8, 1502);
    			attr_dev(main, "class", "pageOverview");
    			add_location(main, file$h, 23, 8, 743);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, t0, anchor);
    			insert_dev(target, main, anchor);
    			append_dev(main, div0);
    			mount_component(panelbalance, div0, null);
    			append_dev(main, t1);
    			append_dev(main, div1);
    			mount_component(panelsend, div1, null);
    			append_dev(main, t2);
    			append_dev(main, div2);
    			mount_component(panelnetworkhashrate, div2, null);
    			append_dev(main, t3);
    			append_dev(main, div3);
    			mount_component(panellocalhashrate, div3, null);
    			append_dev(main, t4);
    			append_dev(main, div4);
    			mount_component(panelstatus, div4, null);
    			append_dev(main, t5);
    			append_dev(main, div5);
    			mount_component(panellatesttx, div5, null);
    			append_dev(main, t6);
    			append_dev(main, div6);
    			append_dev(main, t7);
    			append_dev(main, div7);
    			current = true;
    		},
    		p: noop,
    		i: function intro(local) {
    			if (current) return;
    			transition_in(panelbalance.$$.fragment, local);
    			transition_in(panelsend.$$.fragment, local);
    			transition_in(panelnetworkhashrate.$$.fragment, local);
    			transition_in(panellocalhashrate.$$.fragment, local);
    			transition_in(panelstatus.$$.fragment, local);
    			transition_in(panellatesttx.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(panelbalance.$$.fragment, local);
    			transition_out(panelsend.$$.fragment, local);
    			transition_out(panelnetworkhashrate.$$.fragment, local);
    			transition_out(panellocalhashrate.$$.fragment, local);
    			transition_out(panelstatus.$$.fragment, local);
    			transition_out(panellatesttx.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(t0);
    			if (detaching) detach_dev(main);
    			destroy_component(panelbalance);
    			destroy_component(panelsend);
    			destroy_component(panelnetworkhashrate);
    			destroy_component(panellocalhashrate);
    			destroy_component(panelstatus);
    			destroy_component(panellatesttx);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$h.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class PageOverview extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$h, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PageOverview",
    			options,
    			id: create_fragment$h.name
    		});
    	}
    }

    /* src/components/panels/PanelTxs.svelte generated by Svelte v3.14.0 */

    const file$i = "src/components/panels/PanelTxs.svelte";

    function create_fragment$i(ctx) {
    	let div1;
    	let div0;

    	const block = {
    		c: function create() {
    			div1 = element("div");
    			div0 = element("div");
    			attr_dev(div0, "id", "txs");
    			add_location(div0, file$i, 30, 23, 1044);
    			attr_dev(div1, "class", "rwrap flx");
    			add_location(div1, file$i, 30, 0, 1021);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div1, anchor);
    			append_dev(div1, div0);
    		},
    		p: noop,
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div1);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$i.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    function instance$c($$self, $$props, $$invalidate) {
    	const txs = new Vue({
    			el: "#txs",
    			name: "PanelTransactions",
    			data() {
    				return {
    					pageSettings: {
    						pageSize: 10,
    						pageSizes: [10, 20, 50, 100],
    						pageCount: 5
    					},
    					ddldata: ["All", "generated", "sent", "received", "immature"]
    				};
    			},
    			template: `<div class="rwrap">
	<div class="select-wrap">
		<ejs-dropdownlist id='ddlelement' :dataSource='ddldata' placeholder='Select category to filter'></ejs-dropdownlist>
    </div>
	<ejs-grid :dataSource="this.transactions.txsEx.txs" height="100%" :allowPaging="true" :pageSettings='pageSettings'>
		<e-columns>
			<e-column field='category' headerText='Category' textAlign='Right' width=90></e-column>
			<e-column field='time' headerText='Time' format='unix'  textAlign='Right' width=90></e-column>
			<e-column field='txid' headerText='TxID' textAlign='Right' width=240></e-column>
			<e-column field='amount' headerText='Amount' textAlign='Right' width=90></e-column>
		</e-columns>
	</ejs-grid>
	</div>`,
    			props: { transactions: Object }
    		});

    	$$self.$capture_state = () => {
    		return {};
    	};

    	$$self.$inject_state = $$props => {
    		
    	};

    	return { txs };
    }

    class PanelTxs extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, instance$c, create_fragment$i, safe_not_equal, { txs: 0 });

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PanelTxs",
    			options,
    			id: create_fragment$i.name
    		});
    	}

    	get txs() {
    		return this.$$.ctx.txs;
    	}

    	set txs(value) {
    		throw new Error("<PanelTxs>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}
    }

    /* src/components/pages/PageTransactions.svelte generated by Svelte v3.14.0 */

    function create_fragment$j(ctx) {
    	let t;
    	let current;
    	const paneltxs = new PanelTxs({ $$inline: true });

    	const block = {
    		c: function create() {
    			t = space();
    			create_component(paneltxs.$$.fragment);
    			document.title = "Transactions";
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, t, anchor);
    			mount_component(paneltxs, target, anchor);
    			current = true;
    		},
    		p: noop,
    		i: function intro(local) {
    			if (current) return;
    			transition_in(paneltxs.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(paneltxs.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(t);
    			destroy_component(paneltxs, detaching);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$j.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class PageTransactions extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$j, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PageTransactions",
    			options,
    			id: create_fragment$j.name
    		});
    	}
    }

    /* src/components/pages/PageAddressBook.svelte generated by Svelte v3.14.0 */

    function create_fragment$k(ctx) {
    	let t;

    	const block = {
    		c: function create() {
    			t = text("\n\nPage PageAddressBook");
    			document.title = "Address Book";
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, t, anchor);
    		},
    		p: noop,
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(t);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$k.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class PageAddressBook extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$k, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PageAddressBook",
    			options,
    			id: create_fragment$k.name
    		});
    	}
    }

    /* src/components/pages/PageExplorer.svelte generated by Svelte v3.14.0 */

    function create_fragment$l(ctx) {
    	let t;

    	const block = {
    		c: function create() {
    			t = text("\n\n\nPage PageExplorer");
    			document.title = "Explorer";
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, t, anchor);
    		},
    		p: noop,
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(t);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$l.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class PageExplorer extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$l, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PageExplorer",
    			options,
    			id: create_fragment$l.name
    		});
    	}
    }

    /* src/components/pages/PageSettings.svelte generated by Svelte v3.14.0 */

    const file$j = "src/components/pages/PageSettings.svelte";

    function create_fragment$m(ctx) {
    	let t;
    	let main;

    	const block = {
    		c: function create() {
    			t = space();
    			main = element("main");
    			document.title = "Settings";
    			attr_dev(main, "class", "pageSettings");
    			add_location(main, file$j, 10, 1, 266);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, t, anchor);
    			insert_dev(target, main, anchor);
    		},
    		p: noop,
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(t);
    			if (detaching) detach_dev(main);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$m.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class PageSettings extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$m, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PageSettings",
    			options,
    			id: create_fragment$m.name
    		});
    	}
    }

    /* src/components/pages/PageNotFound.svelte generated by Svelte v3.14.0 */
    const file$k = "src/components/pages/PageNotFound.svelte";

    // (9:0) <Button>
    function create_default_slot$6(ctx) {
    	let t;

    	const block = {
    		c: function create() {
    			t = text("Go Home");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, t, anchor);
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(t);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_default_slot$6.name,
    		type: "slot",
    		source: "(9:0) <Button>",
    		ctx
    	});

    	return block;
    }

    function create_fragment$n(ctx) {
    	let h1;
    	let t1;
    	let t2;
    	let current;

    	const button = new Button({
    			props: {
    				$$slots: { default: [create_default_slot$6] },
    				$$scope: { ctx }
    			},
    			$$inline: true
    		});

    	const block = {
    		c: function create() {
    			h1 = element("h1");
    			h1.textContent = "404 Not Found";
    			t1 = space();
    			create_component(button.$$.fragment);
    			t2 = space();
    			add_location(h1, file$k, 7, 0, 58);
    			document.title = "404";
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, h1, anchor);
    			insert_dev(target, t1, anchor);
    			mount_component(button, target, anchor);
    			insert_dev(target, t2, anchor);
    			current = true;
    		},
    		p: function update(changed, ctx) {
    			const button_changes = {};

    			if (changed.$$scope) {
    				button_changes.$$scope = { changed, ctx };
    			}

    			button.$set(button_changes);
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(button.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(button.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(h1);
    			if (detaching) detach_dev(t1);
    			destroy_component(button, detaching);
    			if (detaching) detach_dev(t2);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$n.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class PageNotFound extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$n, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PageNotFound",
    			options,
    			id: create_fragment$n.name
    		});
    	}
    }

    function createPages() {
        const {subscribe, set} = writable(PageOverview);

        return {
            subscribe,
            overview: () => set(PageOverview),
            transactions: () => set(PageTransactions),
            addressbook: () => set(PageAddressBook),
            explorer: () => set(PageExplorer),
            settings: () => set(PageSettings),
            notfound: () => set(PageNotFound),
        };
    }

    const isPage = createPages();

    /* src/components/ico/Logo.svelte generated by Svelte v3.14.0 */

    const file$l = "src/components/ico/Logo.svelte";

    function create_fragment$o(ctx) {
    	let svg;
    	let path;

    	const block = {
    		c: function create() {
    			svg = svg_element("svg");
    			path = svg_element("path");
    			attr_dev(path, "class", "logofill");
    			attr_dev(path, "d", "M77.08,2.55c3.87,1.03 6.96,2.58 10.32,4.64c5.93,3.87 10.58,8.51 14.19,14.71c3.87,6.19 5.42,13.16 5.42,20.64c0,7.22 -1.81,14.18 -5.41,20.37c-3.61,6.45 -8.25,11.35 -14.19,14.96c-3.35,2.06 -6.96,3.87 -10.32,4.9c-3.87,1.03 -7.74,1.55 -11.61,1.55v-14.45c6.96,-0.26 13.42,-2.58 19.09,-8c5.67,-5.42 8.51,-11.87 8.51,-19.61c0,-7.74 -2.58,-14.19 -7.99,-19.6c-5.42,-5.42 -11.86,-8 -19.6,-8c-7.74,0 -14.44,2.58 -19.6,8c-5.42,5.42 -8,11.87 -8,19.6l0,85.9c-3.1,-3.1 -7.99,-7.74 -13.93,-13.67v-72.23c0,-3.87 0.52,-7.73 1.55,-11.35c1.03,-3.87 2.58,-7.22 4.64,-10.32c3.87,-5.93 8.52,-10.58 14.71,-14.45c6.19,-3.61 13.16,-5.16 20.64,-5.16c3.87,0 8,0.52 11.61,1.55zM78.37,42.28c0,7.22 -5.93,13.16 -13.15,13.16c-7.48,0.26 -13.16,-5.68 -13.16,-13.16c0,-7.22 5.94,-13.16 13.16,-13.16c7.22,0 13.15,5.93 13.15,13.16zM13.63,37.12l0,69.39c-6.19,-6.19 -11.09,-10.83 -13.93,-13.93l0,-55.46z");
    			add_location(path, file$l, 0, 109, 109);
    			attr_dev(svg, "xmlns", "http://www.w3.org/2000/svg");
    			attr_dev(svg, "id", "parallelCoinLogo");
    			attr_dev(svg, "viewBox", "0 0 108 128");
    			attr_dev(svg, "width", "108");
    			attr_dev(svg, "height", "128");
    			attr_dev(svg, "class", "svelte-134e4a6");
    			add_location(svg, file$l, 0, 0, 0);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, svg, anchor);
    			append_dev(svg, path);
    		},
    		p: noop,
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(svg);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$o.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class Logo extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$o, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "Logo",
    			options,
    			id: create_fragment$o.name
    		});
    	}
    }

    /* src/components/layout/Header.svelte generated by Svelte v3.14.0 */

    const file$m = "src/components/layout/Header.svelte";

    function create_fragment$p(ctx) {
    	let header;
    	let div0;
    	let t0;
    	let div1;
    	let t1;
    	let div4;
    	let div3;
    	let div2;
    	let t3;
    	let div5;
    	let h1;
    	let t5;
    	let div6;
    	let t6;
    	let div7;
    	let t7;
    	let div8;
    	let t8;
    	let div9;
    	let button;

    	const block = {
    		c: function create() {
    			header = element("header");
    			div0 = element("div");
    			t0 = space();
    			div1 = element("div");
    			t1 = space();
    			div4 = element("div");
    			div3 = element("div");
    			div2 = element("div");
    			div2.textContent = "ParallelCoin";
    			t3 = space();
    			div5 = element("div");
    			h1 = element("h1");
    			h1.textContent = "Svelte Router";
    			t5 = space();
    			div6 = element("div");
    			t6 = space();
    			div7 = element("div");
    			t7 = space();
    			div8 = element("div");
    			t8 = space();
    			div9 = element("div");
    			button = element("button");
    			button.textContent = "Open";
    			attr_dev(div0, "class", "h1");
    			add_location(div0, file$m, 15, 2, 711);
    			attr_dev(div1, "class", "h2");
    			add_location(div1, file$m, 16, 2, 736);
    			attr_dev(div2, "class", "analysis");
    			add_location(div2, file$m, 19, 4, 813);
    			attr_dev(div3, "class", "searchContent");
    			add_location(div3, file$m, 18, 3, 781);
    			attr_dev(div4, "class", "h3");
    			add_location(div4, file$m, 17, 2, 761);
    			add_location(h1, file$m, 24, 4, 897);
    			attr_dev(div5, "class", "h4");
    			add_location(div5, file$m, 22, 2, 875);
    			attr_dev(div6, "class", "h5");
    			add_location(div6, file$m, 31, 2, 946);
    			attr_dev(div7, "class", "h6");
    			add_location(div7, file$m, 32, 2, 971);
    			attr_dev(div8, "class", "h7");
    			add_location(div8, file$m, 33, 2, 996);
    			attr_dev(button, "id", "toggle");
    			attr_dev(button, "ref", "toggleBoardbtn");
    			attr_dev(button, "class", "e-btn e-info svelte-ryo8nl");
    			attr_dev(button, "cssclass", "e-flat");
    			attr_dev(button, "iconcss", "e-icons burg-icon");
    			attr_dev(button, "istoggle", "true");
    			add_location(button, file$m, 35, 3, 1041);
    			attr_dev(div9, "class", "h8");
    			add_location(div9, file$m, 34, 2, 1021);
    			attr_dev(header, "class", "Header bgLight");
    			add_location(header, file$m, 14, 0, 677);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, header, anchor);
    			append_dev(header, div0);
    			append_dev(header, t0);
    			append_dev(header, div1);
    			append_dev(header, t1);
    			append_dev(header, div4);
    			append_dev(div4, div3);
    			append_dev(div3, div2);
    			append_dev(header, t3);
    			append_dev(header, div5);
    			append_dev(div5, h1);
    			append_dev(header, t5);
    			append_dev(header, div6);
    			append_dev(header, t6);
    			append_dev(header, div7);
    			append_dev(header, t7);
    			append_dev(header, div8);
    			append_dev(header, t8);
    			append_dev(header, div9);
    			append_dev(div9, button);
    		},
    		p: noop,
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(header);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$p.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class Header extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$p, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "Header",
    			options,
    			id: create_fragment$p.name
    		});
    	}
    }

    /* src/components/ico/IcoOverview.svelte generated by Svelte v3.14.0 */

    const file$n = "src/components/ico/IcoOverview.svelte";

    function create_fragment$q(ctx) {
    	let svg;
    	let polygon0;
    	let g;
    	let polygon1;
    	let rect0;
    	let polygon2;
    	let rect1;
    	let rect2;
    	let path;

    	const block = {
    		c: function create() {
    			svg = svg_element("svg");
    			polygon0 = svg_element("polygon");
    			g = svg_element("g");
    			polygon1 = svg_element("polygon");
    			rect0 = svg_element("rect");
    			polygon2 = svg_element("polygon");
    			rect1 = svg_element("rect");
    			rect2 = svg_element("rect");
    			path = svg_element("path");
    			attr_dev(polygon0, "class", "bgMoreLight");
    			attr_dev(polygon0, "points", "42,39 6,39 6,23 24,6 42,23");
    			add_location(polygon0, file$n, 0, 120, 120);
    			attr_dev(polygon1, "points", "39,21 34,16 34,9 39,9");
    			add_location(polygon1, file$n, 0, 204, 204);
    			attr_dev(rect0, "x", "6");
    			attr_dev(rect0, "y", "39");
    			attr_dev(rect0, "width", "36");
    			attr_dev(rect0, "height", "5");
    			add_location(rect0, file$n, 0, 245, 245);
    			attr_dev(g, "class", "bgBlue");
    			add_location(g, file$n, 0, 186, 186);
    			attr_dev(polygon2, "class", "red");
    			attr_dev(polygon2, "points", "24,4.3 4,22.9 6,25.1 24,8.4 42,25.1 44,22.9");
    			add_location(polygon2, file$n, 0, 291, 291);
    			attr_dev(rect1, "x", "18");
    			attr_dev(rect1, "y", "28");
    			attr_dev(rect1, "class", "red");
    			attr_dev(rect1, "width", "12");
    			attr_dev(rect1, "height", "16");
    			add_location(rect1, file$n, 0, 366, 366);
    			attr_dev(rect2, "x", "21");
    			attr_dev(rect2, "y", "17");
    			attr_dev(rect2, "class", "bgBlue");
    			attr_dev(rect2, "width", "6");
    			attr_dev(rect2, "height", "6");
    			add_location(rect2, file$n, 0, 422, 422);
    			attr_dev(path, "class", "bgGreen");
    			attr_dev(path, "d", "M27.5,35.5c-0.3,0-0.5,0.2-0.5,0.5v2c0,0.3,0.2,0.5,0.5,0.5S28,38.3,28,38v-2C28,35.7,27.8,35.5,27.5,35.5z");
    			add_location(path, file$n, 0, 479, 479);
    			attr_dev(svg, "version", "1");
    			attr_dev(svg, "xmlns", "http://www.w3.org/2000/svg");
    			attr_dev(svg, "class", "icon");
    			attr_dev(svg, "viewBox", "0 0 48 48");
    			attr_dev(svg, "enable-background", "new 0 0 48 48");
    			add_location(svg, file$n, 0, 0, 0);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, svg, anchor);
    			append_dev(svg, polygon0);
    			append_dev(svg, g);
    			append_dev(g, polygon1);
    			append_dev(g, rect0);
    			append_dev(svg, polygon2);
    			append_dev(svg, rect1);
    			append_dev(svg, rect2);
    			append_dev(svg, path);
    		},
    		p: noop,
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(svg);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$q.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class IcoOverview extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$q, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "IcoOverview",
    			options,
    			id: create_fragment$q.name
    		});
    	}
    }

    /* src/components/ico/IcoHistory.svelte generated by Svelte v3.14.0 */

    const file$o = "src/components/ico/IcoHistory.svelte";

    function create_fragment$r(ctx) {
    	let svg;
    	let path0;
    	let path1;
    	let g0;
    	let circle0;
    	let circle1;
    	let g1;
    	let path2;
    	let path3;
    	let g2;
    	let rect0;
    	let rect1;
    	let rect2;
    	let rect3;
    	let rect4;
    	let rect5;
    	let rect6;
    	let rect7;
    	let rect8;
    	let rect9;
    	let rect10;
    	let rect11;

    	const block = {
    		c: function create() {
    			svg = svg_element("svg");
    			path0 = svg_element("path");
    			path1 = svg_element("path");
    			g0 = svg_element("g");
    			circle0 = svg_element("circle");
    			circle1 = svg_element("circle");
    			g1 = svg_element("g");
    			path2 = svg_element("path");
    			path3 = svg_element("path");
    			g2 = svg_element("g");
    			rect0 = svg_element("rect");
    			rect1 = svg_element("rect");
    			rect2 = svg_element("rect");
    			rect3 = svg_element("rect");
    			rect4 = svg_element("rect");
    			rect5 = svg_element("rect");
    			rect6 = svg_element("rect");
    			rect7 = svg_element("rect");
    			rect8 = svg_element("rect");
    			rect9 = svg_element("rect");
    			rect10 = svg_element("rect");
    			rect11 = svg_element("rect");
    			attr_dev(path0, "class", "bgMoreLight");
    			attr_dev(path0, "d", "M5,38V14h38v24c0,2.2-1.8,4-4,4H9C6.8,42,5,40.2,5,38z");
    			add_location(path0, file$o, 0, 119, 119);
    			attr_dev(path1, "class", "red");
    			attr_dev(path1, "d", "M43,10v6H5v-6c0-2.2,1.8-4,4-4h30C41.2,6,43,7.8,43,10z");
    			add_location(path1, file$o, 0, 203, 203);
    			attr_dev(circle0, "cx", "33");
    			attr_dev(circle0, "cy", "10");
    			attr_dev(circle0, "r", "3");
    			add_location(circle0, file$o, 0, 303, 303);
    			attr_dev(circle1, "cx", "15");
    			attr_dev(circle1, "cy", "10");
    			attr_dev(circle1, "r", "3");
    			add_location(circle1, file$o, 0, 334, 334);
    			attr_dev(g0, "class", "bgMoreLight");
    			add_location(g0, file$o, 0, 280, 280);
    			attr_dev(path2, "d", "M33,3c-1.1,0-2,0.9-2,2v5c0,1.1,0.9,2,2,2s2-0.9,2-2V5C35,3.9,34.1,3,33,3z");
    			add_location(path2, file$o, 0, 387, 387);
    			attr_dev(path3, "d", "M15,3c-1.1,0-2,0.9-2,2v5c0,1.1,0.9,2,2,2s2-0.9,2-2V5C17,3.9,16.1,3,15,3z");
    			add_location(path3, file$o, 0, 471, 471);
    			attr_dev(g1, "class", "bgGray");
    			add_location(g1, file$o, 0, 369, 369);
    			attr_dev(rect0, "x", "13");
    			attr_dev(rect0, "y", "20");
    			attr_dev(rect0, "width", "4");
    			attr_dev(rect0, "height", "4");
    			add_location(rect0, file$o, 0, 577, 577);
    			attr_dev(rect1, "x", "19");
    			attr_dev(rect1, "y", "20");
    			attr_dev(rect1, "width", "4");
    			attr_dev(rect1, "height", "4");
    			add_location(rect1, file$o, 0, 619, 619);
    			attr_dev(rect2, "x", "25");
    			attr_dev(rect2, "y", "20");
    			attr_dev(rect2, "width", "4");
    			attr_dev(rect2, "height", "4");
    			add_location(rect2, file$o, 0, 661, 661);
    			attr_dev(rect3, "x", "31");
    			attr_dev(rect3, "y", "20");
    			attr_dev(rect3, "width", "4");
    			attr_dev(rect3, "height", "4");
    			add_location(rect3, file$o, 0, 703, 703);
    			attr_dev(rect4, "x", "13");
    			attr_dev(rect4, "y", "26");
    			attr_dev(rect4, "width", "4");
    			attr_dev(rect4, "height", "4");
    			add_location(rect4, file$o, 0, 745, 745);
    			attr_dev(rect5, "x", "19");
    			attr_dev(rect5, "y", "26");
    			attr_dev(rect5, "width", "4");
    			attr_dev(rect5, "height", "4");
    			add_location(rect5, file$o, 0, 787, 787);
    			attr_dev(rect6, "x", "25");
    			attr_dev(rect6, "y", "26");
    			attr_dev(rect6, "width", "4");
    			attr_dev(rect6, "height", "4");
    			add_location(rect6, file$o, 0, 829, 829);
    			attr_dev(rect7, "x", "31");
    			attr_dev(rect7, "y", "26");
    			attr_dev(rect7, "width", "4");
    			attr_dev(rect7, "height", "4");
    			add_location(rect7, file$o, 0, 871, 871);
    			attr_dev(rect8, "x", "13");
    			attr_dev(rect8, "y", "32");
    			attr_dev(rect8, "width", "4");
    			attr_dev(rect8, "height", "4");
    			add_location(rect8, file$o, 0, 913, 913);
    			attr_dev(rect9, "x", "19");
    			attr_dev(rect9, "y", "32");
    			attr_dev(rect9, "width", "4");
    			attr_dev(rect9, "height", "4");
    			add_location(rect9, file$o, 0, 955, 955);
    			attr_dev(rect10, "x", "25");
    			attr_dev(rect10, "y", "32");
    			attr_dev(rect10, "width", "4");
    			attr_dev(rect10, "height", "4");
    			add_location(rect10, file$o, 0, 997, 997);
    			attr_dev(rect11, "x", "31");
    			attr_dev(rect11, "y", "32");
    			attr_dev(rect11, "width", "4");
    			attr_dev(rect11, "height", "4");
    			add_location(rect11, file$o, 0, 1039, 1039);
    			attr_dev(g2, "class", "bgGray");
    			add_location(g2, file$o, 0, 559, 559);
    			attr_dev(svg, "version", "1");
    			attr_dev(svg, "xmlns", "http://www.w3.org/2000/svg");
    			attr_dev(svg, "class", "icon");
    			attr_dev(svg, "viewBox", "0 0 48 48");
    			attr_dev(svg, "enable-background", "new 0 0 48 48");
    			add_location(svg, file$o, 0, 0, 0);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, svg, anchor);
    			append_dev(svg, path0);
    			append_dev(svg, path1);
    			append_dev(svg, g0);
    			append_dev(g0, circle0);
    			append_dev(g0, circle1);
    			append_dev(svg, g1);
    			append_dev(g1, path2);
    			append_dev(g1, path3);
    			append_dev(svg, g2);
    			append_dev(g2, rect0);
    			append_dev(g2, rect1);
    			append_dev(g2, rect2);
    			append_dev(g2, rect3);
    			append_dev(g2, rect4);
    			append_dev(g2, rect5);
    			append_dev(g2, rect6);
    			append_dev(g2, rect7);
    			append_dev(g2, rect8);
    			append_dev(g2, rect9);
    			append_dev(g2, rect10);
    			append_dev(g2, rect11);
    		},
    		p: noop,
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(svg);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$r.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class IcoHistory extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$r, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "IcoHistory",
    			options,
    			id: create_fragment$r.name
    		});
    	}
    }

    /* src/components/ico/IcoAddressBook.svelte generated by Svelte v3.14.0 */

    const file$p = "src/components/ico/IcoAddressBook.svelte";

    function create_fragment$s(ctx) {
    	let svg;
    	let path0;
    	let path1;
    	let path2;

    	const block = {
    		c: function create() {
    			svg = svg_element("svg");
    			path0 = svg_element("path");
    			path1 = svg_element("path");
    			path2 = svg_element("path");
    			attr_dev(path0, "class", "bgLightBlue");
    			attr_dev(path0, "d", "M38,44H12V4h26c2.2,0,4,1.8,4,4v32C42,42.2,40.2,44,38,44z");
    			add_location(path0, file$p, 0, 119, 119);
    			attr_dev(path1, "class", "bgBlue");
    			attr_dev(path1, "d", "M10,4h2v40h-2c-2.2,0-4-1.8-4-4V8C6,5.8,7.8,4,10,4z");
    			add_location(path1, file$p, 0, 207, 207);
    			attr_dev(path2, "class", "bgMoreLight");
    			attr_dev(path2, "d", "M36,24.2c-0.1,4.8-3.1,6.9-5.3,6.7c-0.6-0.1-2.1-0.1-2.9-1.6c-0.8,1-1.8,1.6-3.1,1.6c-2.6,0-3.3-2.5-3.4-3.1 c-0.1-0.7-0.2-1.4-0.1-2.2c0.1-1,1.1-6.5,5.7-6.5c2.2,0,3.5,1.1,3.7,1.3L30,27.2c0,0.3-0.2,1.6,1.1,1.6c2.1,0,2.4-3.9,2.4-4.6 c0.1-1.2,0.3-8.2-7-8.2c-6.9,0-7.9,7.4-8,9.2c-0.5,8.5,6,8.5,7.2,8.5c1.7,0,3.7-0.7,3.9-0.8l0.4,2c-0.3,0.2-2,1.1-4.4,1.1 c-2.2,0-10.1-0.4-9.8-10.8C16.1,23.1,17.4,14,26.6,14C35.8,14,36,22.1,36,24.2z M24.1,25.5c-0.1,1,0,1.8,0.2,2.3 c0.2,0.5,0.6,0.8,1.2,0.8c0.1,0,0.3,0,0.4-0.1c0.2-0.1,0.3-0.1,0.5-0.3c0.2-0.1,0.3-0.3,0.5-0.6c0.2-0.2,0.3-0.6,0.4-1l0.5-5.4 c-0.2-0.1-0.5-0.1-0.7-0.1c-0.5,0-0.9,0.1-1.2,0.3c-0.3,0.2-0.6,0.5-0.9,0.8c-0.2,0.4-0.4,0.8-0.6,1.3S24.2,24.8,24.1,25.5z");
    			add_location(path2, file$p, 0, 284, 284);
    			attr_dev(svg, "version", "1");
    			attr_dev(svg, "xmlns", "http://www.w3.org/2000/svg");
    			attr_dev(svg, "class", "icon");
    			attr_dev(svg, "viewBox", "0 0 48 48");
    			attr_dev(svg, "enable-background", "new 0 0 48 48");
    			add_location(svg, file$p, 0, 0, 0);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, svg, anchor);
    			append_dev(svg, path0);
    			append_dev(svg, path1);
    			append_dev(svg, path2);
    		},
    		p: noop,
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(svg);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$s.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class IcoAddressBook extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$s, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "IcoAddressBook",
    			options,
    			id: create_fragment$s.name
    		});
    	}
    }

    /* src/components/ico/IcoSettings.svelte generated by Svelte v3.14.0 */

    const file$q = "src/components/ico/IcoSettings.svelte";

    function create_fragment$t(ctx) {
    	let svg;
    	let path0;
    	let path1;

    	const block = {
    		c: function create() {
    			svg = svg_element("svg");
    			path0 = svg_element("path");
    			path1 = svg_element("path");
    			attr_dev(path0, "class", "bgGray");
    			attr_dev(path0, "d", "M39.6,27.2c0.1-0.7,0.2-1.4,0.2-2.2s-0.1-1.5-0.2-2.2l4.5-3.2c0.4-0.3,0.6-0.9,0.3-1.4L40,10.8 c-0.3-0.5-0.8-0.7-1.3-0.4l-5,2.3c-1.2-0.9-2.4-1.6-3.8-2.2l-0.5-5.5c-0.1-0.5-0.5-0.9-1-0.9h-8.6c-0.5,0-1,0.4-1,0.9l-0.5,5.5 c-1.4,0.6-2.7,1.3-3.8,2.2l-5-2.3c-0.5-0.2-1.1,0-1.3,0.4l-4.3,7.4c-0.3,0.5-0.1,1.1,0.3,1.4l4.5,3.2c-0.1,0.7-0.2,1.4-0.2,2.2 s0.1,1.5,0.2,2.2L4,30.4c-0.4,0.3-0.6,0.9-0.3,1.4L8,39.2c0.3,0.5,0.8,0.7,1.3,0.4l5-2.3c1.2,0.9,2.4,1.6,3.8,2.2l0.5,5.5 c0.1,0.5,0.5,0.9,1,0.9h8.6c0.5,0,1-0.4,1-0.9l0.5-5.5c1.4-0.6,2.7-1.3,3.8-2.2l5,2.3c0.5,0.2,1.1,0,1.3-0.4l4.3-7.4 c0.3-0.5,0.1-1.1-0.3-1.4L39.6,27.2z M24,35c-5.5,0-10-4.5-10-10c0-5.5,4.5-10,10-10c5.5,0,10,4.5,10,10C34,30.5,29.5,35,24,35z");
    			add_location(path0, file$q, 0, 120, 120);
    			attr_dev(path1, "class", "black");
    			attr_dev(path1, "d", "M24,13c-6.6,0-12,5.4-12,12c0,6.6,5.4,12,12,12s12-5.4,12-12C36,18.4,30.6,13,24,13z M24,30 c-2.8,0-5-2.2-5-5c0-2.8,2.2-5,5-5s5,2.2,5,5C29,27.8,26.8,30,24,30z");
    			add_location(path1, file$q, 0, 839, 839);
    			attr_dev(svg, "version", "1");
    			attr_dev(svg, "xmlns", "http://www.w3.org/2000/svg");
    			attr_dev(svg, "class", "icon");
    			attr_dev(svg, "viewBox", "0 0 48 48");
    			attr_dev(svg, "enable-background", "new 0 0 48 48");
    			add_location(svg, file$q, 0, 0, 0);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, svg, anchor);
    			append_dev(svg, path0);
    			append_dev(svg, path1);
    		},
    		p: noop,
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(svg);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$t.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class IcoSettings extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$t, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "IcoSettings",
    			options,
    			id: create_fragment$t.name
    		});
    	}
    }

    /* src/components/layout/Nav.svelte generated by Svelte v3.14.0 */
    const file$r = "src/components/layout/Nav.svelte";

    function create_fragment$u(ctx) {
    	let nav;
    	let ul;
    	let li0;
    	let button0;
    	let t0;
    	let li1;
    	let button1;
    	let t1;
    	let li2;
    	let button2;
    	let t2;
    	let li3;
    	let button3;
    	let t3;
    	let li4;
    	let button4;
    	let current;
    	let dispose;
    	const icooverview0 = new IcoOverview({ $$inline: true });
    	const icohistory = new IcoHistory({ $$inline: true });
    	const icoaddressbook = new IcoAddressBook({ $$inline: true });
    	const icooverview1 = new IcoOverview({ $$inline: true });
    	const icosettings = new IcoSettings({ $$inline: true });

    	const block = {
    		c: function create() {
    			nav = element("nav");
    			ul = element("ul");
    			li0 = element("li");
    			button0 = element("button");
    			create_component(icooverview0.$$.fragment);
    			t0 = space();
    			li1 = element("li");
    			button1 = element("button");
    			create_component(icohistory.$$.fragment);
    			t1 = space();
    			li2 = element("li");
    			button2 = element("button");
    			create_component(icoaddressbook.$$.fragment);
    			t2 = space();
    			li3 = element("li");
    			button3 = element("button");
    			create_component(icooverview1.$$.fragment);
    			t3 = space();
    			li4 = element("li");
    			button4 = element("button");
    			create_component(icosettings.$$.fragment);
    			attr_dev(button0, "class", "noMargin noPadding noBorder bgTrans sXs cursorPointer");
    			add_location(button0, file$r, 14, 10, 418);
    			attr_dev(li0, "id", "menuoverview");
    			attr_dev(li0, "class", "sidebar-item current");
    			add_location(li0, file$r, 13, 6, 356);
    			attr_dev(button1, "class", "noMargin noPadding noBorder bgTrans sXs cursorPointer");
    			add_location(button1, file$r, 19, 8, 637);
    			attr_dev(li1, "id", "menutransactions");
    			attr_dev(li1, "class", "sidebar-item");
    			add_location(li1, file$r, 18, 6, 581);
    			attr_dev(button2, "class", "noMargin noPadding noBorder bgTrans sXs cursorPointer");
    			add_location(button2, file$r, 24, 8, 856);
    			attr_dev(li2, "id", "menuaddressbook");
    			attr_dev(li2, "class", "sidebar-item");
    			add_location(li2, file$r, 23, 6, 801);
    			attr_dev(button3, "class", "noMargin noPadding noBorder bgTrans sXs cursorPointer");
    			add_location(button3, file$r, 29, 8, 1082);
    			attr_dev(li3, "id", "menublockexplorer");
    			attr_dev(li3, "class", "sidebar-item");
    			add_location(li3, file$r, 28, 6, 1025);
    			attr_dev(button4, "class", "noMargin noPadding noBorder bgTrans sXs cursorPointer");
    			add_location(button4, file$r, 34, 8, 1293);
    			attr_dev(li4, "id", "menusettings");
    			attr_dev(li4, "class", "sidebar-item");
    			add_location(li4, file$r, 33, 6, 1241);
    			attr_dev(ul, "id", "menu");
    			attr_dev(ul, "class", "lsn noPadding");
    			add_location(ul, file$r, 12, 4, 313);
    			attr_dev(nav, "class", "Nav");
    			add_location(nav, file$r, 11, 0, 291);

    			dispose = [
    				listen_dev(button0, "click", isPage.overview, false, false, false),
    				listen_dev(button1, "click", isPage.transactions, false, false, false),
    				listen_dev(button2, "click", isPage.addressbook, false, false, false),
    				listen_dev(button3, "click", isPage.explorer, false, false, false),
    				listen_dev(button4, "click", isPage.settings, false, false, false)
    			];
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, nav, anchor);
    			append_dev(nav, ul);
    			append_dev(ul, li0);
    			append_dev(li0, button0);
    			mount_component(icooverview0, button0, null);
    			append_dev(ul, t0);
    			append_dev(ul, li1);
    			append_dev(li1, button1);
    			mount_component(icohistory, button1, null);
    			append_dev(ul, t1);
    			append_dev(ul, li2);
    			append_dev(li2, button2);
    			mount_component(icoaddressbook, button2, null);
    			append_dev(ul, t2);
    			append_dev(ul, li3);
    			append_dev(li3, button3);
    			mount_component(icooverview1, button3, null);
    			append_dev(ul, t3);
    			append_dev(ul, li4);
    			append_dev(li4, button4);
    			mount_component(icosettings, button4, null);
    			current = true;
    		},
    		p: noop,
    		i: function intro(local) {
    			if (current) return;
    			transition_in(icooverview0.$$.fragment, local);
    			transition_in(icohistory.$$.fragment, local);
    			transition_in(icoaddressbook.$$.fragment, local);
    			transition_in(icooverview1.$$.fragment, local);
    			transition_in(icosettings.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(icooverview0.$$.fragment, local);
    			transition_out(icohistory.$$.fragment, local);
    			transition_out(icoaddressbook.$$.fragment, local);
    			transition_out(icooverview1.$$.fragment, local);
    			transition_out(icosettings.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(nav);
    			destroy_component(icooverview0);
    			destroy_component(icohistory);
    			destroy_component(icoaddressbook);
    			destroy_component(icooverview1);
    			destroy_component(icosettings);
    			run_all(dispose);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$u.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class Nav extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$u, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "Nav",
    			options,
    			id: create_fragment$u.name
    		});
    	}
    }

    /* src/components/layout/UI.svelte generated by Svelte v3.14.0 */
    const file$s = "src/components/layout/UI.svelte";

    function create_fragment$v(ctx) {
    	let t0;
    	let div7;
    	let div6;
    	let div5;
    	let div0;
    	let t1;
    	let t2;
    	let div3;
    	let div1;
    	let t3;
    	let t4;
    	let div2;
    	let t5;
    	let div4;
    	let div4_transition;
    	let current;
    	const logo = new Logo({ $$inline: true });
    	const header = new Header({ $$inline: true });
    	const nav = new Nav({ $$inline: true });
    	var switch_value = ctx.$isPage;

    	function switch_props(ctx) {
    		return { $$inline: true };
    	}

    	if (switch_value) {
    		var switch_instance = new switch_value(switch_props());
    	}

    	const block = {
    		c: function create() {
    			t0 = space();
    			div7 = element("div");
    			div6 = element("div");
    			div5 = element("div");
    			div0 = element("div");
    			create_component(logo.$$.fragment);
    			t1 = space();
    			create_component(header.$$.fragment);
    			t2 = space();
    			div3 = element("div");
    			div1 = element("div");
    			t3 = space();
    			create_component(nav.$$.fragment);
    			t4 = space();
    			div2 = element("div");
    			t5 = space();
    			div4 = element("div");
    			if (switch_instance) create_component(switch_instance.$$.fragment);
    			document.title = "DuOS";
    			attr_dev(div0, "class", "flx fii Logo");
    			add_location(div0, file$s, 30, 5, 513);
    			attr_dev(div1, "class", "Open");
    			add_location(div1, file$s, 33, 6, 611);
    			attr_dev(div2, "class", "Side");
    			add_location(div2, file$s, 35, 6, 656);
    			attr_dev(div3, "class", "Sidebar bgLight");
    			add_location(div3, file$s, 32, 5, 575);
    			attr_dev(div4, "id", "main");
    			attr_dev(div4, "class", "grayGrad Main");
    			add_location(div4, file$s, 38, 0, 694);
    			attr_dev(div5, "class", "grid-container rwrap bgDark");
    			add_location(div5, file$s, 29, 4, 466);
    			attr_dev(div6, "id", "display");
    			attr_dev(div6, "class", "fii");
    			add_location(div6, file$s, 28, 53, 431);
    			attr_dev(div7, "id", "x");
    			attr_dev(div7, "class", "fullScreen bgDark flx lightTheme");
    			add_location(div7, file$s, 28, 0, 378);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, t0, anchor);
    			insert_dev(target, div7, anchor);
    			append_dev(div7, div6);
    			append_dev(div6, div5);
    			append_dev(div5, div0);
    			mount_component(logo, div0, null);
    			append_dev(div5, t1);
    			mount_component(header, div5, null);
    			append_dev(div5, t2);
    			append_dev(div5, div3);
    			append_dev(div3, div1);
    			append_dev(div3, t3);
    			mount_component(nav, div3, null);
    			append_dev(div3, t4);
    			append_dev(div3, div2);
    			append_dev(div5, t5);
    			append_dev(div5, div4);

    			if (switch_instance) {
    				mount_component(switch_instance, div4, null);
    			}

    			current = true;
    		},
    		p: function update(changed, ctx) {
    			if (switch_value !== (switch_value = ctx.$isPage)) {
    				if (switch_instance) {
    					group_outros();
    					const old_component = switch_instance;

    					transition_out(old_component.$$.fragment, 1, 0, () => {
    						destroy_component(old_component, 1);
    					});

    					check_outros();
    				}

    				if (switch_value) {
    					switch_instance = new switch_value(switch_props());
    					create_component(switch_instance.$$.fragment);
    					transition_in(switch_instance.$$.fragment, 1);
    					mount_component(switch_instance, div4, null);
    				} else {
    					switch_instance = null;
    				}
    			}
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(logo.$$.fragment, local);
    			transition_in(header.$$.fragment, local);
    			transition_in(nav.$$.fragment, local);
    			if (switch_instance) transition_in(switch_instance.$$.fragment, local);

    			add_render_callback(() => {
    				if (!div4_transition) div4_transition = create_bidirectional_transition(div4, fade, { duration: 300 }, true);
    				div4_transition.run(1);
    			});

    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(logo.$$.fragment, local);
    			transition_out(header.$$.fragment, local);
    			transition_out(nav.$$.fragment, local);
    			if (switch_instance) transition_out(switch_instance.$$.fragment, local);
    			if (!div4_transition) div4_transition = create_bidirectional_transition(div4, fade, { duration: 300 }, false);
    			div4_transition.run(0);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(t0);
    			if (detaching) detach_dev(div7);
    			destroy_component(logo);
    			destroy_component(header);
    			destroy_component(nav);
    			if (switch_instance) destroy_component(switch_instance);
    			if (detaching && div4_transition) div4_transition.end();
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$v.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    function instance$d($$self, $$props, $$invalidate) {
    	let $isPage;
    	validate_store(isPage, "isPage");
    	component_subscribe($$self, isPage, $$value => $$invalidate("$isPage", $isPage = $$value));
    	let page;

    	const unsubscribe = isPage.subscribe(value => {
    		page = value;
    	});

    	$$self.$capture_state = () => {
    		return {};
    	};

    	$$self.$inject_state = $$props => {
    		if ("page" in $$props) page = $$props.page;
    		if ("$isPage" in $$props) isPage.set($isPage = $$props.$isPage);
    	};

    	return { $isPage };
    }

    class UI extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, instance$d, create_fragment$v, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "UI",
    			options,
    			id: create_fragment$v.name
    		});
    	}
    }

    /* src/DuOS.svelte generated by Svelte v3.14.0 */

    // (9:0) {#if bios.boot}
    function create_if_block_1$4(ctx) {
    	let current;
    	const boot = new Boot({ $$inline: true });

    	const block = {
    		c: function create() {
    			create_component(boot.$$.fragment);
    		},
    		m: function mount(target, anchor) {
    			mount_component(boot, target, anchor);
    			current = true;
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(boot.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(boot.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			destroy_component(boot, detaching);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_if_block_1$4.name,
    		type: "if",
    		source: "(9:0) {#if bios.boot}",
    		ctx
    	});

    	return block;
    }

    // (12:0) {#if !bios.boot}
    function create_if_block$7(ctx) {
    	let current;
    	const ui = new UI({ $$inline: true });

    	const block = {
    		c: function create() {
    			create_component(ui.$$.fragment);
    		},
    		m: function mount(target, anchor) {
    			mount_component(ui, target, anchor);
    			current = true;
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(ui.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(ui.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			destroy_component(ui, detaching);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_if_block$7.name,
    		type: "if",
    		source: "(12:0) {#if !bios.boot}",
    		ctx
    	});

    	return block;
    }

    function create_fragment$w(ctx) {
    	let t;
    	let if_block1_anchor;
    	let current;
    	let if_block0 = bios.boot && create_if_block_1$4(ctx);
    	let if_block1 = !bios.boot && create_if_block$7(ctx);

    	const block = {
    		c: function create() {
    			if (if_block0) if_block0.c();
    			t = space();
    			if (if_block1) if_block1.c();
    			if_block1_anchor = empty();
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			if (if_block0) if_block0.m(target, anchor);
    			insert_dev(target, t, anchor);
    			if (if_block1) if_block1.m(target, anchor);
    			insert_dev(target, if_block1_anchor, anchor);
    			current = true;
    		},
    		p: noop,
    		i: function intro(local) {
    			if (current) return;
    			transition_in(if_block0);
    			transition_in(if_block1);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(if_block0);
    			transition_out(if_block1);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (if_block0) if_block0.d(detaching);
    			if (detaching) detach_dev(t);
    			if (if_block1) if_block1.d(detaching);
    			if (detaching) detach_dev(if_block1_anchor);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$w.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class DuOS extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$w, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "DuOS",
    			options,
    			id: create_fragment$w.name
    		});
    	}
    }

    const duos = new DuOS({
    	target: document.body
    });

    return duos;

}());
//# sourceMappingURL=dui.js.map
